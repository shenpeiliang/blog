# MapLimiter 包

一个 Go 语言实现的并发安全登录尝试限制器，通过跟踪失败的登录尝试并临时阻止可疑 IP 或账户，帮助防止暴力破解攻击。

## 功能特点

- 使用 `sync.Map` 实现线程安全
- 可配置最大尝试次数和阻止时长
- 自动清理过期记录
- 简单的 API 用于记录失败尝试和检查访问权限
- 基于时间的阻止机制，支持自定义时间格式

## 安装

```go
import "cms/util/limiter"
```

## 使用方式

### 基础用法

```go
// 创建带有默认配置的限流器
tf := timex.NewTimeFormatter(timex.DefaultTimeLayout)
limiter := limiter.NewMapLimiter(tf)

// 记录一次失败尝试
limiter.RecordFailedAttempt("user@example.com")

// 检查是否允许登录
if limiter.Allowed("user@example.com") {
    // 允许登录
} else {
    // 阻止登录
}

// 重置某个key的尝试次数
limiter.Reset("user@example.com")
```

### 自定义配置

```go
// 创建带有自定义配置的限流器
limiter := limiter.NewMapLimiter(
    timex.NewTimeFormatter(timex.DefaultTimeLayout),
    limiter.WithMapMaxAttempts(3),               // 3次尝试后阻止
    limiter.WithMapBlockDuration(10*time.Minute), // 阻止10分钟
)
```

## API 参考

### 类型

#### `LoginAttempt`
```go
type LoginAttempt struct {
    Count     uint      // 失败尝试次数
    LastTry   time.Time // 最后一次尝试时间
    BlockedAt time.Time // 被阻止的时间
}
```

#### `MapLimiter`
```go
type MapLimiter struct {
    // 包含未导出字段
}
```

### 方法

#### `NewMapLimiter`
```go
func NewMapLimiter(timeFormatter *timex.TimeFormatter, options ...MapLimiterOption) *MapLimiter
```
创建一个新的 MapLimiter 实例，可传入可选配置。

配置选项:
- `WithMapMaxAttempts(maxAttempts uint)` - 设置阻止前的最大尝试次数
- `WithMapBlockDuration(blockDuration time.Duration)` - 设置阻止时长

#### `RecordFailedAttempt`
```go
func (l *MapLimiter) RecordFailedAttempt(key string)
```
记录指定 key (如用户名或 IP) 的一次失败登录尝试。

#### `Allowed`
```go
func (l *MapLimiter) Allowed(key string) bool
```
检查指定 key 是否允许登录。如果因尝试次数过多被阻止则返回 false。

#### `Reset`
```go
func (l *MapLimiter) Reset(key string)
```
重置指定 key 的尝试计数。

#### `GetAttempt`
```go
func (l *MapLimiter) GetAttempt(key string) (*LoginAttempt, bool)
```
获取指定 key 的登录尝试记录(如果存在)。

#### `CleanupExpired`
```go
func (l *MapLimiter) CleanupExpired()
```
清理过期的尝试记录(由清理协程自动调用)。

#### `Stop`
```go
func (l *MapLimiter) Stop()
```
停止后台清理协程。

## 默认值

- 默认最大尝试次数: 5 次
- 默认阻止时长: 5 分钟

## 清理机制

限流器会自动运行一个清理协程，移除超过两倍阻止时长的记录。这有助于防止积累陈旧记录导致内存泄漏。

## 完整代码

```go
package limiter

import (
	"cms/util/timex"
	"sync"
	"time"
)

const (
	DefaultMaxAttempts   = 5
	DefaultBlockDuration = 5 * time.Minute
)

type LoginAttempt struct {
	Count     uint
	LastTry   time.Time
	BlockedAt time.Time
}

type MapLimiterOptions struct {
	MaxAttempts   uint
	BlockDuration time.Duration
}

type MapLimiterOption func(*MapLimiterOptions)

// 设置最大尝试次数
func WithMapMaxAttempts(maxAttempts uint) MapLimiterOption {
	return func(o *MapLimiterOptions) { o.MaxAttempts = maxAttempts }
}

// 设置阻止时间
func WithMapBlockDuration(blockDuration time.Duration) MapLimiterOption {
	return func(o *MapLimiterOptions) { o.BlockDuration = blockDuration }
}

type MapLimiter struct {
	TimeFormatter *timex.TimeFormatter
	loginAttempts sync.Map
	maxAttempts   uint
	blockDuration time.Duration
	stopChan      chan struct{}
}

// 创建MapLimiter实例
func NewMapLimiter(timeFormatter *timex.TimeFormatter, options ...MapLimiterOption) *MapLimiter {
	opts := MapLimiterOptions{
		MaxAttempts:   DefaultMaxAttempts,
		BlockDuration: DefaultBlockDuration,
	}

	for _, option := range options {
		option(&opts)
	}

	limiter := &MapLimiter{
		TimeFormatter: timeFormatter,
		maxAttempts:   opts.MaxAttempts,
		blockDuration: opts.BlockDuration,
		stopChan:      make(chan struct{}),
	}

	// 启动定期清理
	go limiter.startCleanupRoutine()

	return limiter
}

// 启动清理协程
func (l *MapLimiter) startCleanupRoutine() {
	ticker := time.NewTicker(l.blockDuration)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.CleanupExpired()
		case <-l.stopChan:
			return
		}
	}
}

// 停止清理协程
func (l *MapLimiter) Stop() {
	close(l.stopChan)
}

// 记录错误尝试次数
func (l *MapLimiter) RecordFailedAttempt(key string) {
	now := l.TimeFormatter.Now()

	if attempt, exists := l.loginAttempts.Load(key); exists {
		att := attempt.(*LoginAttempt)
		att.Count++
		att.LastTry = now

		if att.Count >= l.maxAttempts {
			att.BlockedAt = now
		}
	} else {
		l.loginAttempts.Store(key, &LoginAttempt{
			Count:   1,
			LastTry: now,
		})
	}
}

// 检查是否允许登录
func (l *MapLimiter) Allowed(key string) bool {
	if attempt, exists := l.loginAttempts.Load(key); exists {
		att := attempt.(*LoginAttempt)

		if att.Count >= l.maxAttempts {
			if time.Since(att.BlockedAt) < l.blockDuration {
				return false
			}
			att.Count = 0
		}
	}
	return true
}

// 清理过期的尝试记录
func (l *MapLimiter) CleanupExpired() {
	expirationTime := l.blockDuration * 2

	l.loginAttempts.Range(func(key, value interface{}) bool {
		att := value.(*LoginAttempt)
		if time.Since(att.LastTry) > expirationTime {
			l.loginAttempts.Delete(key)
		}
		return true
	})
}

// 重置尝试次数
func (l *MapLimiter) Reset(key string) {
	l.loginAttempts.Delete(key)
}

// 获取尝试记录
func (l *MapLimiter) GetAttempt(key string) (*LoginAttempt, bool) {
	if attempt, exists := l.loginAttempts.Load(key); exists {
		return attempt.(*LoginAttempt), true
	}
	return nil, false
}

```