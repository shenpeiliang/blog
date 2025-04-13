# MultipartRequest 文档

## 概述

`MultipartRequest` 是一个用于构建和发送 multipart/form-data 请求的 Go 语言封装库，支持文件上传和表单字段提交。

## 功能特性

- 支持 POST 和 GET 请求方法
- 简单易用的链式调用 API
- 支持文件上传和多种类型表单字段
- 完整的响应处理，包括状态码、头部和响应体
- 自动处理 multipart 表单构建

## 核心结构

### MultipartRequest

```go
type MultipartRequest struct {
    URL    string            // 请求URL
    Fields map[string]any    // 表单字段
    Files  map[string]string // 文件字段
    writer *multipart.Writer 
    buffer *bytes.Buffer
}
```

### Response

```go
type Response struct {
    StatusCode int               // HTTP状态码
    Headers    map[string]string // 响应头
    Body       []byte            // 原始响应体
    Result     string            `json:"result"`  // 解析后的结果
    Version    string            `json:"version"` // 解析后的版本
}
```

## 使用方法

### 初始化

```go
req := NewMultipartRequest()
```

### 添加表单字段

```go
req.AddField("username", "john")
   .AddField("age", 25)
```

### 添加文件

```go
req.AddFile("avatar", "/path/to/avatar.jpg")
```

### 发送请求

```go
// POST 请求
resp, err := req.Send(http.MethodPost, "https://example.com/upload")

// GET 请求
resp, err := req.Send(http.MethodGet, "https://example.com/api")
```

## 完整示例

### 文件上传示例

```go
func uploadFile() error {
    req := NewMultipartRequest()
    resp, err := req.
        AddField("description", "profile picture").
        AddFile("file", "/path/to/photo.jpg").
        Send(http.MethodPost, "https://example.com/upload")
    
    if err != nil {
        return fmt.Errorf("上传失败: %v", err)
    }
    
    fmt.Printf("上传成功! 状态码: %d\n", resp.StatusCode)
    fmt.Printf("响应结果: %s\n", resp.Result)
    return nil
}
```

### GET 请求示例

```go
func getData() error {
    resp, err := NewMultipartRequest().
        AddField("page", 1).
        AddField("limit", 20).
        Send(http.MethodGet, "https://example.com/api/data")
    
    if err != nil {
        return err
    }
    
    fmt.Printf("获取到数据: %s\n", resp.Body)
    return nil
}
```
## 完整代码

```go
package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// 封装了multipart/form-data请求的构建和发送
type MultipartRequest struct {
	URL    string            // 请求的URL
	Fields map[string]any    // 表单字段，支持任意类型的值
	Files  map[string]string // 文件字段，key为字段名，value为文件路径
	writer *multipart.Writer // multipart.Writer
	buffer *bytes.Buffer     // 请求体缓冲区
}

type Response struct {
	StatusCode int               // HTTP 状态码
	Headers    map[string]string // 响应头
	Body       []byte            // 原始响应体
	Result     string            `json:"result"`  // 响应结果
	Version    string            `json:"version"` // 版本号
}

// 创建一个新的MultipartRequest实例
func NewMultipartRequest() *MultipartRequest {
	return &MultipartRequest{
		URL:    "",
		Fields: make(map[string]any),
		Files:  make(map[string]string),
		buffer: &bytes.Buffer{},
	}
}

// 添加一个表单字段，支持任意类型的值
func (m *MultipartRequest) AddField(name string, value any) *MultipartRequest {
	m.Fields[name] = value
	return m
}

// 添加一个文件字段
func (m *MultipartRequest) AddFile(fieldName, filePath string) *MultipartRequest {
	m.Files[fieldName] = filePath
	return m
}

// 构建multipart请求体
func (m *MultipartRequest) buildPost() error {
	// 重新初始化buffer和writer
	m.buffer.Reset()
	m.writer = multipart.NewWriter(m.buffer)

	// 添加表单字段
	for key, value := range m.Fields {
		if err := m.writer.WriteField(key, fmt.Sprintf("%v", value)); err != nil {
			return fmt.Errorf("构建表单字段 %s 参数出错 %v", key, err)
		}
	}

	// 添加文件字段
	for fieldName, filePath := range m.Files {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("构建表单字段 %s 参数出错 %v", filePath, err)
		}

		//defer file.Close() // 在循环中 defer 不会立即执行

		part, err := m.writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			file.Close()
			return fmt.Errorf("构建表单字段 %s 参数出错 %v", fieldName, err)
		}

		if _, err := io.Copy(part, file); err != nil {
			file.Close()
			return fmt.Errorf("文件操作出错: %v", err)
		}
	}

	// 关闭writer
	if err := m.writer.Close(); err != nil {
		return fmt.Errorf("文件操作出错: %v", err)
	}

	return nil
}

// 构建GET请求
func (m *MultipartRequest) buildGet() {
	// 对于GET请求，将表单字段附加到URL中
	values := url.Values{}
	for key, value := range m.Fields {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	m.URL += "?" + values.Encode()
}

// Send 发送请求，支持GET和POST方法
func (m *MultipartRequest) Send(method, url string) (Response, error) {
	m.URL = url

	var response Response

	// 构建请求体
	var body io.Reader
	var contentType string

	if method == http.MethodPost {
		if err := m.buildPost(); err != nil {
			return response, err
		}
		body = m.buffer
		contentType = m.writer.FormDataContentType()
	} else {
		m.buildGet()
		body = nil
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, m.URL, body)
	if err != nil {
		return response, fmt.Errorf("请求创建出错: %v", err)
	}

	// 设置请求头
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("请求发送出错: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应头
	response.StatusCode = resp.StatusCode
	response.Headers = make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			response.Headers[k] = v[0]
		}
	}

	// 读取响应
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("读取响应数据出错: %v", err)
	}
	response.Body = bodyData

	// 解析响应
	if err := json.Unmarshal(bodyData, &response); err != nil {
		response.Result = string(bodyData)
		return response, fmt.Errorf("解析响应数据出错: %v", err)
	}

	return response, nil
}

```