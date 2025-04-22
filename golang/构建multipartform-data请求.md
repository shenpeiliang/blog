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
package multipart

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
	"time"
)

// Request 封装了multipart/form-data请求的构建和发送
type Request struct {
	URL     string            // 请求的URL
	Fields  map[string]any    // 表单字段，支持任意类型的值
	Files   map[string]string // 文件字段，key为字段名，value为文件路径
	Headers map[string]string // 自定义请求头
	Client  *http.Client      // HTTP客户端
	writer  *multipart.Writer // multipart.Writer
	buffer  *bytes.Buffer     // 请求体缓冲区
	timeout time.Duration     // 请求超时时间
}

// Response 表示HTTP响应
type Response struct {
	StatusCode int               // HTTP 状态码
	Headers    map[string]string // 响应头
	Body       []byte            // 原始响应体
	Result     string            `json:"result"`  // 响应结果
	Version    string            `json:"version"` // 版本号
}

// NewRequest 创建一个新的Request实例
func NewRequest(opts ...Option) *Request {
	req := &Request{
		Fields:  make(map[string]any),
		Files:   make(map[string]string),
		Headers: make(map[string]string),
		buffer:  &bytes.Buffer{},
		Client:  &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// Option 配置选项类型
type Option func(*Request)

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(m *Request) {
		m.timeout = timeout
		m.Client.Timeout = timeout
	}
}

// WithClient 设置自定义HTTP客户端
func WithClient(client *http.Client) Option {
	return func(m *Request) {
		m.Client = client
	}
}

// WithHeader 设置请求头
func WithHeader(key, value string) Option {
	return func(m *Request) {
		m.Headers[key] = value
	}
}

// AddField 添加一个表单字段，支持任意类型的值
func (m *Request) AddField(name string, value any) *Request {
	m.Fields[name] = value
	return m
}

// AddFile 添加一个文件字段
func (m *Request) AddFile(fieldName, filePath string) *Request {
	m.Files[fieldName] = filePath
	return m
}

// AddHeader 添加请求头
func (m *Request) AddHeader(key, value string) *Request {
	m.Headers[key] = value
	return m
}

// buildPost 构建multipart请求体
func (m *Request) buildPost() error {
	m.buffer.Reset()
	m.writer = multipart.NewWriter(m.buffer)

	// 添加表单字段
	for key, value := range m.Fields {
		if err := m.writer.WriteField(key, fmt.Sprintf("%v", value)); err != nil {
			return fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// 添加文件字段
	for fieldName, filePath := range m.Files {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		part, err := m.writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to create form file %s: %w", fieldName, err)
		}

		if _, err := io.Copy(part, file); err != nil {
			file.Close()
			return fmt.Errorf("failed to copy file content: %w", err)
		}
		file.Close()
	}

	return m.writer.Close()
}

// buildGet 构建GET请求
func (m *Request) buildGet() error {
	values := url.Values{}
	for key, value := range m.Fields {
		values.Add(key, fmt.Sprintf("%v", value))
	}

	u, err := url.Parse(m.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	query := u.Query()
	for k, v := range values {
		query[k] = v
	}
	u.RawQuery = query.Encode()
	m.URL = u.String()

	return nil
}

// Send 发送请求，支持GET和POST方法
func (m *Request) Send(method, url string) (Response, error) {
	var response Response
	m.URL = url

	// 验证HTTP方法
	if method != http.MethodGet && method != http.MethodPost {
		return response, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// 构建请求
	var body io.Reader
	var contentType string

	if method == http.MethodPost {
		if err := m.buildPost(); err != nil {
			return response, fmt.Errorf("failed to build POST request: %w", err)
		}
		body = m.buffer
		contentType = m.writer.FormDataContentType()
	} else {
		if err := m.buildGet(); err != nil {
			return response, fmt.Errorf("failed to build GET request: %w", err)
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, m.URL, body)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range m.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := m.Client.Do(req)
	if err != nil {
		return response, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 处理响应
	response.StatusCode = resp.StatusCode
	response.Headers = make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			response.Headers[k] = v[0]
		}
	}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}
	response.Body = bodyData

	// 尝试解析JSON响应
	if err := json.Unmarshal(bodyData, &response); err != nil {
		response.Result = string(bodyData)
	}

	return response, nil
}

```