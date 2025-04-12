# AuthLimiter 限流器包文档

## 概述

AuthLimiter 是一个基于 Redis 的认证限流器实现，用于防止暴力破解攻击。它通过记录失败尝试次数并在达到阈值时锁定账户来实现保护机制。

## 功能特性

- 基于 Redis 实现，支持分布式环境
- 可配置的尝试次数限制和锁定时间
- 滑动窗口计数机制
- 原子性操作保证线程安全
- 简单的 API 接口

## 安装

```go
import "cms/util/limiter"
```

## 快速开始

### 基本用法

```go
// 初始化 Redis 缓存
redisCache := cache.NewRedisCache(...)

// 创建限流器实例
limiter := limiter.NewAuthLimiter(redisCache)

// 检查是否允许尝试
if limiter.Allow("user123") {
    // 处理认证逻辑
    if authFailed {
        // 认证失败时增加计数
        err := limiter.IncrLimitCount("user123")
        if err != nil {
            // 处理错误
        }
    }
} else {
    // 账户被锁定
    remaining, _ := limiter.GetRemainingTime("user123")
    fmt.Printf("账户已锁定，剩余时间: %v", remaining)
}
```

### 自定义配置

```go
limiter := limiter.NewAuthLimiter(redisCache,
    limiter.WithMaxAttempts(3),               // 最多尝试3次
    limiter.WithLockDuration(1*time.Hour),    // 锁定1小时
    limiter.WithWindowDuration(10*time.Minute), // 10分钟窗口
    limiter.WithKeyPrefix("myapp:auth:"),     // 自定义键前缀
)
```


## API 文档

### `NewAuthLimiter`

```go
func NewAuthLimiter(redisCache *cache.RedisCache, options ...AuthLimiterOption) *AuthLimiter
```

创建新的限流器实例。

参数:
- `redisCache`: Redis 缓存实例
- `options`: 可选的配置函数

### `Allow`

```go
func (s *AuthLimiter) Allow(key string) bool
```

检查指定键是否允许尝试。如果返回 false 表示已被锁定。

### `IncrLimitCount`

```go
func (s *AuthLimiter) IncrLimitCount(key string) error
```

增加指定键的尝试计数。当达到最大尝试次数时会自动锁定。

### `GetRemainingTime`

```go
func (s *AuthLimiter) GetRemainingTime(key string) (time.Duration, error)
```

获取指定键的锁定剩余时间。

## 配置选项

| 选项 | 默认值 | 描述 |
|------|--------|------|
| WithMaxAttempts | 5 | 最大尝试次数 |
| WithLockDuration | 30分钟 | 锁定持续时间 |
| WithWindowDuration | 5分钟 | 计数窗口时间 |
| WithKeyPrefix | "limiter:auth:" | Redis 键前缀 |

## 实现细节

1. 使用 Redis 存储尝试计数和锁定状态
2. 使用 Lua 脚本保证原子性操作
3. 采用滑动窗口机制进行计数
4. 锁定后自动设置过期时间

## 最佳实践

1. 对于关键账户系统，建议设置较低的尝试次数(如3次)
2. 生产环境建议使用自定义键前缀避免冲突
3. 锁定时间应根据业务需求合理设置
4. 窗口时间应大于正常用户的认证间隔

## 完整代码
```go
package limiter

import (
	"cms/service/cache"
	"cms/util/errorx"
	"time"
)

// 默认配置常量
const (
	defaultMaxAttempts    = 5
	defaultLockDuration   = 30 * time.Minute
	defaultWindowDuration = 5 * time.Minute
	defaultKeyPrefix      = "limiter:auth:"
)

type AuthLimiter struct {
	redisCache     *cache.RedisCache
	maxAttempts    uint
	lockDuration   time.Duration
	windowDuration time.Duration
	attemptPrefix  string
	lockPrefix     string
}

type AuthLimiterOptions struct {
	maxAttempts    uint
	lockDuration   time.Duration
	windowDuration time.Duration
	keyPrefix      string
}

// 配置选项函数类型
type AuthLimiterOption func(*AuthLimiterOptions)

// 设置最大尝试次数
func WithMaxAttempts(n uint) AuthLimiterOption {
	return func(opts *AuthLimiterOptions) {
		opts.maxAttempts = n
	}
}

// 设置锁定持续时间
func WithLockDuration(d time.Duration) AuthLimiterOption {
	return func(opts *AuthLimiterOptions) {
		opts.lockDuration = d
	}
}

// 设置计数窗口时间
func WithWindowDuration(d time.Duration) AuthLimiterOption {
	return func(opts *AuthLimiterOptions) {
		opts.windowDuration = d
	}
}

// 设置键前缀
func WithKeyPrefix(prefix string) AuthLimiterOption {
	return func(opts *AuthLimiterOptions) {
		opts.keyPrefix = prefix
	}
}

// 创建新的限流器实例
//
// redis缓存实例
//
// opts 配置选项，可通过 WithMaxAttempts、WithLockDuration、WithWindowDuration、WithKeyPrefix 等函数设置
func NewAuthLimiter(redisCache *cache.RedisCache, options ...AuthLimiterOption) *AuthLimiter {
	// 设置默认值
	opts := AuthLimiterOptions{
		maxAttempts:    defaultMaxAttempts,
		lockDuration:   defaultLockDuration,
		windowDuration: defaultWindowDuration,
		keyPrefix:      defaultKeyPrefix,
	}

	// 应用选项
	for _, option := range options {
		option(&opts)
	}

	return &AuthLimiter{
		redisCache:     redisCache,
		maxAttempts:    opts.maxAttempts,
		lockDuration:   opts.lockDuration,
		windowDuration: opts.windowDuration,
		attemptPrefix:  opts.keyPrefix + ":attempt:",
		lockPrefix:     opts.keyPrefix + ":lock:",
	}
}

// 检查是否允许访问
func (s *AuthLimiter) Allow(key string) bool {
	lockKey := s.lockPrefix + key
	exists, err := s.redisCache.Exists(lockKey)
	return err == nil && !exists
}

// 增加尝试计数
func (s *AuthLimiter) IncrLimitCount(key string) error {
	attemptKey := s.attemptPrefix + key
	lockKey := s.lockPrefix + key

	luaScript := `
	local attemptKey = KEYS[1]
	local lockKey = KEYS[2]
	local maxAttempts = tonumber(ARGV[1])
	local lockDuration = tonumber(ARGV[2])
	local windowDuration = tonumber(ARGV[3])
	
	if redis.call("EXISTS", lockKey) == 1 then
		return 0
	end
	
	local count = redis.call("INCR", attemptKey)
	
	if count == 1 then
		redis.call("EXPIRE", attemptKey, windowDuration)
	end
	
	if count >= maxAttempts then
		redis.call("SET", lockKey, 1, "EX", lockDuration)
	end
	
	return count
	`

	_, err := s.redisCache.EvalCmd(luaScript, []string{attemptKey, lockKey},
		s.maxAttempts,
		int(s.lockDuration.Seconds()),
		int(s.windowDuration.Seconds())).Int64()

	if err != nil {
		return errorx.Wrap(err, "系统错误，请稍后再试")
	}

	return nil
}

// 获取锁定剩余时间
func (s *AuthLimiter) GetRemainingTime(key string) (time.Duration, error) {
	seconds, err := s.redisCache.TTL(s.lockPrefix + key)
	return time.Duration(seconds) * time.Second, err
}

```