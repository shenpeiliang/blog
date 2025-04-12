# Redis 分布式锁实现文档  

## 1. 功能概述  
提供基于 Redis 的分布式锁机制，支持：  
- **获取锁**（`Lock`）：支持重试机制，可自定义重试次数或无限重试。  
- **释放锁**（`Unlock`）：删除锁键，释放资源。  

---

## 2. 核心方法  

### 2.1 `Lock` - 获取分布式锁  
**方法签名**：  
```go
func (rc *RedisCache) Lock(
    key string,          // 锁的唯一标识
    value any,           // 锁的值（通常为唯一标识，如UUID）
    timeout time.Duration, // 锁的自动超时时间
    retryNum ...uint8,   // 可选参数：最大重试次数（默认无限重试）
) (isLocked bool, err error)
```

**参数说明**：  
| 参数       | 类型             | 说明                                                                 |
|------------|------------------|----------------------------------------------------------------------|
| `key`      | `string`         | 锁的键名，实际存储时会添加前缀 `CACHE_PREFIX_LOCK`。                |
| `value`    | `any`            | 锁的值（建议使用唯一标识，避免误删其他客户端的锁）。                |
| `timeout`  | `time.Duration`  | 锁的自动过期时间（防止死锁）。                                      |
| `retryNum` | `...uint8`       | 可选参数：最大重试次数。未指定时默认无限重试，指定时按次数重试。    |

**返回值**：  
- `isLocked`：是否成功获取锁。  
- `err`：错误信息（如 Redis 连接失败或重试次数耗尽）。  

**逻辑流程**：  
1. 使用 `SETNX` 尝试设置锁键（仅当键不存在时成功）。  
2. 若成功，立即返回 `true`。  
3. 若失败：  
   - 如果指定了 `retryNum` 且剩余次数 ≤ 0，返回错误。  
   - 否则等待 100ms 后重试。  

---

### 2.2 `Unlock` - 释放分布式锁  
**方法签名**：  
```go
func (rc *RedisCache) Unlock(key string) (err error)
```

**参数说明**：  
| 参数    | 类型     | 说明                                |
|---------|----------|-------------------------------------|
| `key`   | `string` | 锁的键名（与 `Lock` 的 `key` 一致）。 |

**返回值**：  
- `err`：删除锁键时的错误（如 Redis 连接失败）。  

**注意事项**：  
- 需确保 `value` 的唯一性，避免误删其他客户端的锁（需结合 Lua 脚本实现原子性检查）。  

---

## 3. 使用示例  

### 3.1 获取锁（无限重试）  
```go
lockKey := "order_123"
lockValue := uuid.New().String()
timeout := 10 * time.Second

isLocked, err := redisCache.Lock(lockKey, lockValue, timeout)
if err != nil {
    log.Fatal("获取锁失败:", err)
}
defer redisCache.Unlock(lockKey)

if isLocked {
    fmt.Println("执行业务逻辑...")
}
```

### 3.2 获取锁（指定重试次数）  
```go
isLocked, err := redisCache.Lock(lockKey, lockValue, timeout, 3) // 最多重试3次
if err != nil {
    log.Fatal("获取锁失败:", err)
}
```

## 完整代码

```go
package cache

import (
	"cms/util/errorx"
	"cms/util/setting"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var (

	//缓存默认时间
	defaultCacheTimeout = 86400 * time.Second

	///缓存默认超时
	defaultContextTimeout = 5 * time.Second

	//缓存读取默认错误
	DefaultReadError = "系统缓存出错"
)

type RedisCache struct {
	Client *redis.Client
}

// 初始化Redis连接
func NewRedisClient() (*RedisCache, error) {
	db, _ := strconv.Atoi(setting.Viper.GetString("redis.database"))
	timeout := time.Duration(setting.Viper.GetInt("redis.timeout")) * time.Second

	rc := redis.NewClient(&redis.Options{
		Network:      setting.Viper.GetString("redis.network"),
		Addr:         setting.Viper.GetString("redis.addr"),
		Username:     setting.Viper.GetString("redis.username"),
		Password:     setting.Viper.GetString("redis.password"),
		DB:           db,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})

	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, errorx.Wrap(err, "无法连接到 Redis")
	}

	return &RedisCache{Client: rc}, nil
}

// 关闭Redis连接
func (rc *RedisCache) Close() {
	if rc.Client != nil {
		rc.Client.Close()
	}
}

// 设置分布式锁
//
// retryNum: 重试次数，默认一直重试,如果指定了重试次数，则到达重试次数时还没有获得锁，则返回失败
//
// 返回值：
// isLocked: 是否成功获取锁
// err: 错误信息
func (rc *RedisCache) Lock(key string, value any, timeout time.Duration, retryNum ...uint8) (isLocked bool, err error) {
	var (
		remaining   uint8 = 0     // 默认重试次数为0, 即一直重试
		useRetryNum bool  = false // 是否使用指定的重试次数，默认一直重试
	)

	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	// 如果提供了重试次数参数，则使用该值
	if len(retryNum) > 0 && retryNum[0] > 0 {
		useRetryNum = true
		remaining = retryNum[0]
	}

	for {
		isLocked, err = rc.Client.SetNX(ctx, CACHE_PREFIX_LOCK+key, value, timeout).Result()
		if err != nil {
			err = errorx.Wrap(err, DefaultReadError)
			return
		}

		//设置成功，跳出循环
		if isLocked {
			return true, nil
		}

		//重试次数超出，返回失败
		if useRetryNum {
			if remaining <= 0 {
				err = errorx.New("系统提示：获取操作权限失败，请稍后再试")
				return
			}

			// 减少重试次数
			remaining--
		}

		// 等待100毫秒，继续尝试获取锁
		time.Sleep(100 * time.Millisecond)
	}

}

// 释放分布式锁
func (rc *RedisCache) Unlock(key string) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	err = rc.Client.Del(ctx, CACHE_PREFIX_LOCK+key).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	return nil
}

// 设置过期时间
func (rc *RedisCache) Expire(key string, expiration time.Duration) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	err = rc.Client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	return nil
}

// 获取事务管道
func (rc *RedisCache) TxPipeline() redis.Pipeliner {
	return rc.Client.TxPipeline()
}

// 累加计数器
func (rc *RedisCache) Incr(key string) (newValue int64, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	newValue, err = rc.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errorx.Wrap(err, DefaultReadError)
	}

	return
}

// 获取超时上下文
func getContextWithTimeout(t ...time.Duration) (context.Context, context.CancelFunc) {
	timeout := defaultContextTimeout
	if len(t) > 0 {
		timeout = t[0]
	}

	return context.WithTimeout(context.Background(), timeout)
}

// 设置缓存
func (rc *RedisCache) Set(key string, value any, timeout ...time.Duration) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	expiration := defaultCacheTimeout
	if len(timeout) > 0 {
		expiration = timeout[0]
	}

	err = rc.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	return nil
}

// 设置缓存
func (rc *RedisCache) SetJson(key string, value any, timeout ...time.Duration) (err error) {
	data, err := sonic.Marshal(value)
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	expiration := defaultCacheTimeout
	if len(timeout) > 0 {
		expiration = timeout[0]
	}

	return rc.Set(key, data, expiration)
}

// 获取缓存
func (rc *RedisCache) Get(key string) (data string, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	data, err = rc.Client.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		err = errorx.Wrap(err, DefaultReadError)
	}

	return
}

// 获取缓存
func (rc *RedisCache) GetJson(key string, result any) (err error) {
	data, err := rc.Get(key)
	if err != nil {
		return
	}

	if data == "" {
		return
	}

	err = sonic.UnmarshalString(data, &result)
	if err != nil {
		err = errorx.Wrap(err, DefaultReadError)
		return
	}

	return
}

// 清除缓存
func (rc *RedisCache) Clear(key ...string) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	err = rc.Client.Del(ctx, key...).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	return nil
}

// 追加集合
func (rc *RedisCache) SAdd(key string, member any, timeout ...time.Duration) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	expiration := defaultCacheTimeout
	if len(timeout) > 0 {
		expiration = timeout[0]
	}

	err = rc.Client.SAdd(ctx, key, member).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	if expiration > 0 {
		err = rc.Client.Expire(ctx, key, expiration).Err()
		if err != nil {
			return errorx.Wrap(err, DefaultReadError)
		}
	}

	return nil
}

// 判断成员元素是否是集合的成员
func (rc *RedisCache) SIsMember(key string, member any) (is bool, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	is, err = rc.Client.SIsMember(ctx, key, member).Result()
	if err != nil {
		err = errorx.Wrap(err, DefaultReadError)
		return
	}

	return
}

// 追加或更新有序集合成员
func (rc *RedisCache) ZAdd(key string, members ...redis.Z) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	//不设置过期时间，默认是永不过期
	err = rc.Client.ZAdd(ctx, key, members...).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	return nil
}

// 有序集合分数最高的成员
//
// fmt.Printf("分数值最高的成员： %s, 分数： %f", scores[0].Member, scores[0].Score)
func (rc *RedisCache) ZRevRangeWithScores(key string, start, stop int64) (scores []redis.Z, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	scores, err = rc.Client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		err = errorx.Wrap(err, DefaultReadError)
		return
	}

	return
}

// 哈希存储
//
// struct类型
func (rc *RedisCache) HSet(key string, value any, timeout ...time.Duration) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	expiration := defaultCacheTimeout
	if len(timeout) > 0 {
		expiration = timeout[0]
	}

	err = rc.Client.HSet(ctx, key, value).Err()
	if err != nil {
		return errorx.Wrap(err, DefaultReadError)
	}

	if expiration > 0 {
		err = rc.Client.Expire(ctx, key, expiration).Err()
		if err != nil {
			return errorx.Wrap(err, DefaultReadError)
		}
	}

	return nil
}

// 哈希获取值
//
// struct类型
func (rc *RedisCache) HGetAll(key string, dest any) (err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	err = rc.Client.HGetAll(ctx, key).Scan(dest)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return errorx.Wrap(err, DefaultReadError)
	}

	return
}

// 获取缓存命令
func (rc *RedisCache) GetCmd(key string) *redis.StringCmd {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	return rc.Client.Get(ctx, key)
}

// 设置缓存命令
func (rc *RedisCache) EvalCmd(script string, keys []string, args ...any) *redis.Cmd {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	return rc.Client.Eval(ctx, script, keys, args)
}

// 检查key是否存在
func (rc *RedisCache) Exists(key string) (exists bool, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	result, err := rc.Client.Exists(ctx, key).Result()
	if err != nil {
		err = errorx.Wrap(err, DefaultReadError)
		return
	}

	exists = result > 0

	return
}

// 获取key的过期时间
func (rc *RedisCache) TTL(key string) (ttl time.Duration, err error) {
	//设置超时
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	ttl, err = rc.Client.TTL(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		err = errorx.Wrap(err, DefaultReadError)
		return
	}

	return
}

```