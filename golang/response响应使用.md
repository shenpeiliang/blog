# HTTP 响应工具包文档

## 概述

本工具包提供了一套标准的 HTTP 响应处理函数，用于在 Iris Web 框架中统一返回 JSON 格式的响应数据。包含成功响应、失败响应和错误处理等功能。

## 功能特性

- **统一响应格式**：标准化 JSON 响应结构
- **错误处理**：自动区分自定义错误和标准错误
- **日志记录**：自动记录错误日志
- **高性能**：使用 sonic 进行 JSON 序列化
- **类型安全**：严格的参数类型检查

## 安装

```go
import "your-module-path/util"
```

## 核心函数

### 1. 成功响应

#### Success

```go
func Success(ctx iris.Context, msg string, data ...any)
```

- **参数**：
  - `ctx`：Iris 上下文
  - `msg`：成功消息
  - `data`：可选，返回的数据内容

- **响应格式**：
```json
{
  "code": "200",
  "msg": "操作成功",
  "data": {...}
}
```

### 2. 失败响应

#### Fail

```go
func Fail(ctx iris.Context, msg string, code ...string)
```

- **参数**：
  - `ctx`：Iris 上下文
  - `msg`：错误消息
  - `code`：可选，错误码（默认 "-1"）

- **响应格式**：
```json
{
  "code": "-1",
  "msg": "操作失败",
  "data": {}
}
```

#### FailError

```go
func FailError(ctx iris.Context, err error)
```

- **参数**：
  - `ctx`：Iris 上下文
  - `err`：错误对象

- **功能**：
  - 自动识别自定义错误类型 `errorx.CustomError`
  - 自动记录错误日志
  - 返回标准化错误响应

### 3. 内部函数

#### response

```go
func response(ctx iris.Context, code, msg string, data any)
```

- **参数**：
  - `ctx`：Iris 上下文
  - `code`：状态码
  - `msg`：消息
  - `data`：数据内容

- **功能**：
  - 核心响应处理函数
  - 设置 JSON Content-Type
  - 处理 JSON 序列化错误

#### DebugError

```go
func DebugError(ctx iris.Context, err error)
```

- **参数**：
  - `ctx`：Iris 上下文
  - `err`：错误对象

- **功能**：
  - 记录错误日志
  - 从上下文中获取 logger 实例

## 使用示例

### 基本使用

```go
package main

import (
	"your-module-path/util"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	app.Get("/success", func(ctx iris.Context) {
		data := map[string]string{
			"name": "John",
			"age":  "30",
		}
		util.Success(ctx, "操作成功", data)
	})

	app.Get("/fail", func(ctx iris.Context) {
		util.Fail(ctx, "操作失败", "400")
	})

	app.Get("/error", func(ctx iris.Context) {
		err := errors.New("数据库连接失败")
		util.FailError(ctx, err)
	})

	app.Listen(":8080")
}
```

### 自定义错误处理

```go
app.Get("/custom-error", func(ctx iris.Context) {
    err := errorx.NewCustomError(1001, "参数验证失败")
    util.FailError(ctx, err)
    
    // 响应格式：
    // {
    //   "code": "1001",
    //   "msg": "参数验证失败",
    //   "data": {}
    // }
})
```

## 响应格式说明

所有响应都遵循以下 JSON 格式：

```json
{
  "code": "string",  // 状态码
  "msg": "string",   // 消息
  "data": any        // 数据内容
}
```

- **成功响应**：code = "200"
- **失败响应**：code ≠ "200"

## 完整代码

```go
package util

import (
	"cms/util/errorx"
	"cms/util/logx"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/kataras/iris/v12"
)

// 返回信息
func response(ctx iris.Context, code, msg string, data any) {
	content := map[string]any{
		"code": code,
		"msg":  msg,
		"data": data,
	}

	// 设置响应头 Content-Type 为 application/json
	ctx.ContentType("application/json")

	contentJson, err := sonic.Marshal(content)
	if err != nil {
		content["code"] = "-1"
		content["msg"] = "系统出错，请稍后重试"
		content["data"] = map[string]any{}

		ctx.JSON(content)
		return
	}

	_, err = ctx.Write(contentJson)

	//记录错误日志
	DebugError(ctx, err)
}

// 成功返回
func Success(ctx iris.Context, msg string, data ...any) {
	var d any
	if len(data) > 0 {
		d = data[0]
	} else {
		d = map[string]any{}
	}

	response(ctx, "200", msg, d)
}

// 失败返回
func Fail(ctx iris.Context, msg string, code ...string) {
	errCode := "-1"
	if len(code) > 0 {
		errCode = code[0]
	}
	response(ctx, errCode, msg, map[string]any{})
}

// 失败返回
func FailError(ctx iris.Context, err error) {
	customError := new(errorx.CustomError)

	code := "-1"
	var message string

	//是否是自定义的error
	if errors.As(err, &customError) {
		code = strconv.Itoa(customError.Code)
		message = customError.ErrorMessage()
	} else {
		//标准error
		message = err.Error()
	}

	//记录错误日志
	DebugError(ctx, err)

	response(ctx, code, message, map[string]any{})
}

// 记录错误日志
func DebugError(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	logger, ok := ctx.Values().Get("logger").(*logx.Logger)
	if !ok {
		fmt.Println("Error: logger not found or type assertion failed")
		return
	}

	logger.Log.Error(err)
}

```