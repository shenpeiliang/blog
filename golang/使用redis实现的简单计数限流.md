# CountLimiter 限流器包文档

## 概述

CountLimiter 是一个基于 Redis 的简单计数器限流器实现，用于限制在指定时间窗口内的请求数量。它采用固定窗口算法实现基础的限流功能。

## 功能特性

- 基于 Redis 实现，支持分布式环境
- 固定时间窗口限流算法
- 可配置的请求限制数和时间窗口
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
limiter := limiter.NewCountLimiter(redisCache)

// 检查是否允许请求
allowed, err := limiter.Allow("api:user:list")
if err != nil {
    // 处理错误
}
if allowed {
    // 处理请求
} else {
    // 请求被限流
}
```

### 自定义配置

```go
limiter := limiter.NewCountLimiter(redisCache,
    limiter.WithLimit(500),            // 每分钟最多500次请求
    limiter.WithWindow(30*time.Second), // 30秒时间窗口
)
```

## 完整代码

```go
package limiter

import (
	"cms/service/cache"
	"cms/util/errorx"
	"time"
)

const (
	defaultLimit  = 100             // 默认请求限制数
	defaultWindow = 1 * time.Minute // 默认时间窗口
)

type CountLimiter struct {
	cache  *cache.RedisCache
	limit  int64         // 限制的总请求数量
	window time.Duration // 限制的时间窗口
}

type CountLimiterOptions struct {
	limit  int64
	window time.Duration
}

type CountLimiterOption func(*CountLimiterOptions)

func WithLimit(limit int64) CountLimiterOption {
	return func(opts *CountLimiterOptions) {
		opts.limit = limit
	}
}

func WithWindow(window time.Duration) CountLimiterOption {
	return func(opts *CountLimiterOptions) {
		opts.window = window
	}
}

func NewCountLimiter(cache *cache.RedisCache, options ...CountLimiterOption) *CountLimiter {
	opts := CountLimiterOptions{
		limit:  defaultLimit,
		window: defaultWindow,
	}

	for _, option := range options {
		option(&opts)
	}

	return &CountLimiter{
		cache:  cache,
		limit:  opts.limit,
		window: opts.window,
	}
}

func (l *CountLimiter) Allow(key string) (bool, error) {
	script := `
	local current = redis.call('incr', KEYS[1])
	if current == 1 then
		redis.call('expire', KEYS[1], ARGV[1])
	end
	return current
	`

	currentCount, err := l.cache.EvalCmd(script, []string{key}, l.window.Seconds()).Int64()
	if err != nil {
		return false, errorx.Wrap(err, "execute rate limit script failed")
	}

	return currentCount <= l.limit, nil
}

func (l *CountLimiter) GetLimit() int64 {
	return l.limit
}

func (l *CountLimiter) GetWindow() time.Duration {
	return l.window
}
```

## API 文档

### `NewCountLimiter`

```go
func NewCountLimiter(cache *cache.RedisCache, options ...CountLimiterOption) *CountLimiter
```

创建新的计数器限流器实例。

参数:
- `cache`: Redis 缓存实例
- `options`: 可选的配置函数

### `Allow`

```go
func (l *CountLimiter) Allow(key string) (bool, error)
```

检查指定键是否允许请求。返回:
- `true`: 允许请求
- `false`: 请求被限流
- `error`: 执行过程中发生的错误

### `GetLimit`

```go
func (l *CountLimiter) GetLimit() int64
```

获取当前配置的请求限制数。

### `GetWindow`

```go
func (l *CountLimiter) GetWindow() time.Duration
```

获取当前配置的时间窗口。

## 配置选项

| 选项 | 默认值 | 描述 |
|------|--------|------|
| WithLimit | 100 | 时间窗口内的最大请求数 |
| WithWindow | 1分钟 | 限流时间窗口 |

## 实现细节

1. 使用 Redis INCR 命令实现原子计数器
2. 首次计数时设置过期时间实现固定窗口
3. 使用 Lua 脚本保证原子性操作
4. 简单高效的固定窗口算法

## 最佳实践

1. 根据业务需求合理设置限制数和时间窗口
2. 对于不同的接口使用不同的 key 前缀
3. 监控限流触发情况以调整限流参数
4. 结合其他限流算法实现更精细的控制

## 注意事项

1. 固定窗口算法在窗口边界可能出现请求突增
2. 需要确保 Redis 的高可用性
3. 时间窗口不宜设置过短，避免频繁操作 Redis
