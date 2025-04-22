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
	"context"
	"io"
	"net/http"
	"time"
)

const (
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
)

// 请求结果
type Result struct {
	Body  []byte
	Error error
}

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

// GetAsync 发起异步GET请求，返回结果通道
func (r *Requester) GetAsync(url string, header map[string]string) <-chan Result {
	resultChan := make(chan Result, 1)
	go r.doAsyncRequest(RequestMethodGet, url, header, nil, resultChan)
	return resultChan
}

// PostAsync 发起异步POST请求，返回结果通道
func (r *Requester) PostAsync(url string, header map[string]string, jsonBody []byte) <-chan Result {
	resultChan := make(chan Result, 1)
	go r.doAsyncRequest(RequestMethodPost, url, header, jsonBody, resultChan)
	return resultChan
}

// 执行异步HTTP请求
func (r *Requester) doAsyncRequest(method, url string, header map[string]string, body []byte, resultChan chan<- Result) {
	body, err := r.doRequest(method, url, header, body)
	resultChan <- Result{Body: body, Error: err}
	close(resultChan)
}

// Get 发起 GET 请求
func (r *Requester) Get(url string, header map[string]string) ([]byte, error) {
	return r.doRequest(RequestMethodGet, url, header, nil)
}

// Post 发起 POST 请求
func (r *Requester) Post(url string, header map[string]string, jsonBody []byte) ([]byte, error) {
	return r.doRequest(RequestMethodPost, url, header, jsonBody)
}

// 执行HTTP请求
func (r *Requester) doRequest(method, url string, header map[string]string, body []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

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

	req = req.WithContext(ctx)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

```