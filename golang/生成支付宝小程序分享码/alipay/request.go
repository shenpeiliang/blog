package alipay

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

func WithAlipayPublicKeyPath(path string) RequesterOption {
	return func(r *Requester) {
		r.alipayPublicKeyPath = path
	}
}

// Requester 结构体
type Requester struct {
	client              *http.Client
	timeout             time.Duration
	alipayPublicKeyPath string
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

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//签名验证
	err = verifySignByCert(resp, result, r.alipayPublicKeyPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}
