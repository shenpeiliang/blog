# RandomService 模块文档

## 概述

`RandomService` 是一个提供随机字符串生成功能的 Go 模块，采用接口抽象和选项模式设计，支持灵活的随机数生成策略。

## 功能特性

- **可插拔的随机数生成策略**：通过 `RandomMethod` 接口支持不同实现
- **选项模式配置**：支持初始化时配置或运行时修改生成策略
- **默认实现**：内置 UUID 生成器作为默认实现

## 接口定义

### RandomMethod 接口

```go
type RandomMethod interface {
    Get() string
}
```

任何实现了 `Get() string` 方法的类型都可以作为随机数生成策略。

## 核心结构体

### RandomService

```go
type RandomService struct {
    Method RandomMethod
}
```

- `Method`：当前使用的随机数生成策略

## 使用方法

### 1. 初始化服务

#### 使用默认实现 (UUID)

```go
service := NewRandomService()
```

#### 使用自定义实现

```go
service := NewRandomService(WithRandomMethod(&MyRandomGenerator{}))
```

### 2. 获取随机字符串

```go
randomStr := service.GetRandChars()
```

### 3. 运行时更改生成策略

```go
service.SetRandomMethod(&MyRandomGenerator{})
```

## 选项函数

### WithRandomMethod

```go
func WithRandomMethod(m RandomMethod) Option
```

用于在初始化时指定随机数生成策略。

## 示例代码

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    // 使用默认UUID生成器
    service := NewRandomService()
    fmt.Println(service.GetRandChars()) // 输出类似: "5f4dcc3b-5aa9-4b2c-9b9a-8e9b6b2a7c1d"
    
    // 使用自定义生成器
    service.SetRandomMethod(&SimpleRandom{})
    fmt.Println(service.GetRandChars()) // 输出类似: "abc123"
}
```

### 完整代码

```go
package token

// 接口
type RandomMethod interface {
	Get() string
}

type RandomService struct {
	Method RandomMethod
}

// 选项函数类型
type Option func(*RandomService)

// 设置随机方法的选项函数
func WithRandomMethod(m RandomMethod) Option {
	return func(s *RandomService) {
		s.Method = m
	}
}

// 初始化
func NewRandomService(opts ...Option) *RandomService {
	service := &RandomService{
		Method: NewUuidService(), // 默认使用uuid生成随机字符串
	}

	// 应用选项
	for _, opt := range opts {
		opt(service)
	}

	return service
}

// 设置随机方法
func (s *RandomService) SetRandomMethod(m RandomMethod) *RandomService {
	s.Method = m
	return s
}

// 获取随机字符串
func (s *RandomService) GetRandChars() string {
	return s.Method.Get()
}

```

### 随机字符串生成器

```go
package token

import "math/rand"

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type RandStringService struct {
	Number int
	Chars  []rune
}

func NewRandStringService(number int, chars []rune) *RandStringService {
	return &RandStringService{
		Number: number,
		Chars:  chars,
	}
}

// 随机字符串
func (s *RandStringService) Get() string {

	var letters []rune

	if len(s.Chars) == 0 {
		letters = defaultLetters
	} else {
		letters = s.Chars
	}

	b := make([]rune, s.Number)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

```

### uuid生成器

```go
package token

import uuid "github.com/satori/go.uuid"

type UuidService struct {
}

func NewUuidService() *UuidService {
	return &UuidService{}
}

// 随机字符串
func (*UuidService) Get() string {
	return uuid.NewV4().String()
}

```

### jwt生成器

```go
package token

import (
	"cms/util/setting"
	"cms/util/timex"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

type JwtService struct {
	TimeFormatter *timex.TimeFormatter
}

func NewJwtService(timeFormatter *timex.TimeFormatter) *JwtService {
	return &JwtService{
		TimeFormatter: timeFormatter,
	}
}

// 获取存储的对象
func GetJwtClaims(ctx iris.Context) jwt.MapClaims {

	value := ctx.Values().Get("jwt").(*jwt.Token)

	return value.Claims.(jwt.MapClaims)
}

// 随机字符串
func (s *JwtService) Get() string {
	//保留的claims 有固定的含义和用途
	claims := jwt.MapClaims{
		"iss": setting.Viper.GetString("jwt.iss"),                            //发行者
		"exp": s.TimeFormatter.UnixTime() + setting.Viper.GetUint("jwt.exp"), //过期时间
		"sub": setting.Viper.GetString("jwt.sub"),                            //主题
		"aud": setting.Viper.GetString("jwt.aud"),                            //用户
	}

	//创建JWT
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, claims)

	// 使用签名算法对JWT进行签名
	t, _ := token.SignedString([]byte(setting.Viper.GetString("jwt.secret")))

	return t
}

```