# HTTP请求工具使用指南

## 概述

`request` 包提供了一个简单易用的HTTP客户端封装，支持GET和POST请求，具有超时控制和请求头设置功能，支持使用通信(channel)来协调异步请求。


## 快速开始

### 基本使用

```go
// 创建默认Requester实例（30秒超时）
requester := request.NewRequester()

// GET请求
body, err := requester.Get("https://example.com", nil)
if err != nil {
    // 处理错误
}
fmt.Println(string(body))

// POST请求
postData := []byte(`{"key":"value"}`)
headers := map[string]string{
    "Content-Type": "application/json",
}
resp, err := requester.Post("https://example.com/api", headers, postData)
if err != nil {
    // 处理错误
}
fmt.Println(string(resp))
```

## 高级配置

### 自定义配置

```go
// 使用自定义配置
requester := request.NewRequester(
    request.WithTimeout(10*time.Second), // 设置10秒超时
    request.WithClient(&http.Client{     // 自定义HTTP客户端
        Transport: &http.Transport{
            MaxIdleConns: 10,
        },
    }),
)
```

### 配置选项

| 选项 | 描述 | 示例 |
|------|------|------|
| `WithTimeout` | 设置请求超时时间 | `WithTimeout(10*time.Second)` |
| `WithClient` | 使用自定义HTTP客户端 | `WithClient(&http.Client{...})` |

## API参考

### `NewRequester`

创建新的Requester实例

```go
func NewRequester(opts ...RequesterOption) *Requester
```

**参数**:
- `opts`: 可选配置项

**返回值**:
- `*Requester`: Requester实例

### `Get`

发起GET请求

```go
func (r *Requester) Get(url string, header map[string]string) ([]byte, error)
```

**参数**:
- `url`: 请求URL
- `header`: 请求头(map[string]string)

**返回值**:
- `[]byte`: 响应体
- `error`: 错误信息

### `Post`

发起POST请求

```go
func (r *Requester) Post(url string, header map[string]string, jsonBody []byte) ([]byte, error)
```

**参数**:
- `url`: 请求URL
- `header`: 请求头(map[string]string)
- `jsonBody`: POST请求体([]byte)

**返回值**:
- `[]byte`: 响应体
- `error`: 错误信息

## 最佳实践

1. **重用Requester实例**:
   ```go
   // 在应用初始化时创建
   var requester = request.NewRequester()
   
   // 在需要的地方重复使用
   func handler() {
       body, _ := requester.Get("https://example.com", nil)
       // ...
   }
   ```

2. **错误处理**:
   ```go
   body, err := requester.Get(url, headers)
   if err != nil {
       if errors.Is(err, context.DeadlineExceeded) {
           // 处理超时
       } else {
           // 处理其他错误
       }
   }
   ```

3. **设置合理的超时**:
   ```go
   // 根据请求类型设置不同超时
   apiRequester := request.NewRequester(
       request.WithTimeout(5*time.Second),
   )
   
   downloadRequester := request.NewRequester(
       request.WithTimeout(30*time.Second),
   )
   ```

## 示例

### 带认证头的API请求

```go
requester := request.NewRequester()

headers := map[string]string{
    "Content-Type":  "application/json",
    "Authorization": "Bearer xxxxx",
}

data, err := requester.Get("https://api.example.com/data", headers)
if err != nil {
    // 处理错误
}

// 处理响应数据
```

### 提交JSON数据

```go
requester := request.NewRequester()

payload := []byte(`{"name":"John","age":30}`)
headers := map[string]string{
    "Content-Type": "application/json",
}

response, err := requester.Post("https://api.example.com/users", headers, payload)
if err != nil {
    // 处理错误
}

// 处理响应
```

## 完整代码
```go
package request

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const (
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
)

type RequesterOption func(*Requester)

func WithClient(client *http.Client) RequesterOption {
	return func(r *Requester) {
		r.client = client
	}
}

func WithTimeout(timeout time.Duration) RequesterOption {
	return func(r *Requester) {
		r.timeout = timeout
	}
}

// Requester 结构体
type Requester struct {
	client  *http.Client
	timeout time.Duration
}

// NewRequester 创建新的 Requester 实例
func NewRequester(opts ...RequesterOption) *Requester {
	r := &Requester{
		client:  &http.Client{},
		timeout: 30 * time.Second, // 默认超时
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Get 发起 GET 请求（同步）
func (r *Requester) Get(url string, header map[string]string) ([]byte, error) {
	return r.doRequest(RequestMethodGet, url, header, nil)
}

// Post 发起 POST 请求（同步）
func (r *Requester) Post(url string, header map[string]string, jsonBody []byte) ([]byte, error) {
	return r.doRequest(RequestMethodPost, url, header, jsonBody)
}

// 执行HTTP请求
func (r *Requester) doRequest(method, url string, header map[string]string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

```

支持异步
```go
package async

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout    = 30 * time.Second
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
)

// 异步请求结果
type Result struct {
	Body  []byte
	Error error
}

// 请求器配置选项
type RequesterOption func(*Requester)

// 设置http client
func WithClient(client *http.Client) RequesterOption {
	return func(r *Requester) {
		r.client = client
	}
}

// 设置超时时间
func WithTimeout(timeout time.Duration) RequesterOption {
	return func(r *Requester) {
		r.client.Timeout = timeout
	}
}

type Requester struct {
	client *http.Client
}

// NewRequester 创建一个请求器
func NewRequester(opts ...RequesterOption) *Requester {
	r := &Requester{
		client: &http.Client{
			Timeout: defaultTimeout, //设置默认超时时间
		},
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Requester) Get(ctx context.Context, url string, header map[string]string) ([]byte, error) {
	return r.doRequest(ctx, RequestMethodGet, url, header, nil)
}

func (r *Requester) Post(ctx context.Context, url string, header map[string]string, jsonBody []byte) ([]byte, error) {
	return r.doRequest(ctx, RequestMethodPost, url, header, jsonBody)
}

func (r *Requester) GetAsync(ctx context.Context, url string, header map[string]string, resultChan chan<- Result) {
	go r.asyncWrapper(ctx, RequestMethodGet, url, header, nil, resultChan)
}

func (r *Requester) PostAsync(ctx context.Context, url string, header map[string]string, jsonBody []byte, resultChan chan<- Result) {
	go r.asyncWrapper(ctx, RequestMethodPost, url, header, jsonBody, resultChan)
}

// 异步请求的包装器
func (r *Requester) asyncWrapper(ctx context.Context, method, url string, header map[string]string, body []byte, resultChan chan<- Result) {
	defer func() {
		if err := recover(); err != nil {
			select {
			case <-ctx.Done():
				// 如果上下文被取消，直接退出
				return
			case resultChan <- Result{Error: fmt.Errorf("panic: %v", err)}:
			}
		}
	}()

	// 请求前检查上下文是否已经被取消
	select {
	case <-ctx.Done():
		// 如果上下文被取消，直接退出
		return
	default:
		// 继续执行
	}

	// 执行请求
	body, err := r.doRequest(ctx, method, url, header, body)

	// 再次检查上下文是否已经被取消
	select {
	case <-ctx.Done():
		// 如果上下文被取消，直接退出
		return
	case resultChan <- Result{Body: body, Error: err}:
	}
}

// 发送请求
func (r *Requester) doRequest(ctx context.Context, method, url string, header map[string]string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	return data, nil
}

```