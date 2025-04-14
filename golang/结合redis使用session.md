# Session Manager 使用文档

## 概述

本文档介绍了如何使用基于Redis的Session管理器。该管理器使用Iris框架的session组件，并将会话数据存储在Redis中。

## 功能特性

- 基于Redis的会话存储
- 可配置的会话过期时间
- 线程安全
- 简单的API接口

## 安装与初始化

首先确保项目中已添加相关依赖：

```go
import (
	"github.com/kataras/iris/v12/sessions"
	"cms/util/session"
)
```

初始化Session管理器：

```go
sessionManager := session.NewSessionManager()
```

## 集成到Iris应用

```go
app := iris.New()
sessionManager := session.NewSessionManager()
app.Use(sessionManager.Handler())
```

## 基本使用

### 设置会话值

```go
func setHandler(ctx iris.Context) {
    sess := sessions.Get(ctx)
    sess.Set("username", "john_doe")
    ctx.WriteString("Session value set")
}
```

### 获取会话值

```go
func getHandler(ctx iris.Context) {
    sess := sessions.Get(ctx)
    username := sess.GetString("username")
    ctx.Writef("Username: %s", username)
}
```

### 删除会话值

```go
func deleteHandler(ctx iris.Context) {
    sess := sessions.Get(ctx)
    sess.Delete("username")
    ctx.WriteString("Session value deleted")
}
```

### 销毁整个会话

```go
func destroyHandler(ctx iris.Context) {
    sess := sessions.Get(ctx)
    sess.Destroy()
    ctx.WriteString("Session destroyed")
}
```

## 配置选项

Session管理器支持通过配置文件(`setting.Viper`)自定义Redis连接参数：

```yaml
redis:
  network: "tcp"
  addr: "localhost:6379"
  timeout: 30
  max_active: 10
  password: ""
  database: "0"
  prefix: "myapp_"
  username: ""
```

配置项说明：
- `network`: 网络类型(默认tcp)
- `addr`: Redis服务器地址
- `timeout`: 连接超时(秒)
- `max_active`: 最大活动连接数
- `password`: Redis密码
- `database`: 数据库编号
- `prefix`: 键前缀
- `username`: 用户名(Redis 6.0+)

## 会话过期

默认配置为浏览器会话结束时过期(`Expires: -1`)。如需修改过期时间，可在`NewSessionManager()`函数中调整：

```go
Expires: 24 * time.Hour // 24小时后过期
```

## 示例代码

完整示例参考Iris官方文档：
[https://github.com/kataras/iris/blob/main/_examples/sessions/basic/main.go](https://github.com/kataras/iris/blob/main/_examples/sessions/basic/main.go)

## 完整代码

```go
package session

import (
	"cms/util/setting"
	"time"

	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

type SessionManager struct {
	Manager *sessions.Sessions
}

// 初始化session管理器
func NewSessionManager() *SessionManager {
	// 初始化session管理器
	manager := sessions.New(sessions.Config{
		Cookie:  sessions.DefaultCookieName,
		Expires: -1, //浏览器关闭后session过期
	})

	manager.UseDatabase(redis.New(loadRedisConfig()))

	return &SessionManager{
		Manager: manager,
	}
}

// 加载redis配置
func loadRedisConfig() redis.Config {
	// 默认配置
	config := redis.DefaultConfig()

	// 是否有配置项
	if network := setting.Viper.GetString("redis.network"); network != "" {
		config.Network = network
		config.Addr = setting.Viper.GetString("redis.addr")
		config.Timeout = time.Duration(setting.Viper.GetInt("redis.timeout")) * time.Second
		config.MaxActive = setting.Viper.GetInt("redis.max_active")
		config.Password = setting.Viper.GetString("redis.password")
		config.Database = setting.Viper.GetString("redis.database")
		config.Prefix = setting.Viper.GetString("redis.prefix")
		config.Username = setting.Viper.GetString("redis.username")
	}

	return config
}

```