package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestBuilder struct {
	client  *Client
	method  string
	url     string
	headers http.Header
	queries map[string][]string
	body    io.Reader
	err     error
}

func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.method = strings.ToUpper(method)
	return rb
}

func (rb *RequestBuilder) Get(url string) *RequestBuilder {
	rb.method = http.MethodGet
	rb.url = url
	return rb
}

func (rb *RequestBuilder) Post(url string) *RequestBuilder {
	rb.method = http.MethodPost
	rb.url = url
	return rb
}

func (rb *RequestBuilder) Put(url string) *RequestBuilder {
	rb.method = http.MethodPut
	rb.url = url
	return rb
}

func (rb *RequestBuilder) Delete(url string) *RequestBuilder {
	rb.method = http.MethodDelete
	rb.url = url
	return rb
}

func (rb *RequestBuilder) URL(u string) *RequestBuilder {
	rb.url = u
	return rb
}

func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	rb.headers.Set(key, value)
	return rb
}

func (rb *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		rb.headers.Set(k, v)
	}
	return rb
}

func (rb *RequestBuilder) Query(key, value string) *RequestBuilder {
	rb.queries[key] = append(rb.queries[key], value)
	return rb
}

func (rb *RequestBuilder) QueryParams(params map[string]string) *RequestBuilder {
	for k, v := range params {
		rb.Query(k, v)
	}
	return rb
}

func (rb *RequestBuilder) Body(body io.Reader) *RequestBuilder {
	rb.body = body
	return rb
}

func (rb *RequestBuilder) JSONBody(v any) *RequestBuilder {
	data, err := json.Marshal(v)
	if err != nil {
		rb.err = fmt.Errorf("marshal json body: %w", err)
		return rb
	}
	rb.body = bytes.NewReader(data)
	rb.headers.Set("Content-Type", "application/json")
	return rb
}

func (rb *RequestBuilder) Do(ctx context.Context) *ResponseHandler {
	if rb.err != nil {
		return &ResponseHandler{err: rb.err}
	}

	// 构建完整 URL
	fullURL := rb.url
	if rb.client.baseURL != "" {
		fullURL = rb.client.baseURL + rb.url
	}

	// 添加查询参数
	if len(rb.queries) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return &ResponseHandler{err: fmt.Errorf("parse url: %w", err)}
		}
		q := u.Query()
		for k, vals := range rb.queries {
			for _, v := range vals {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, rb.method, fullURL, rb.body)
	if err != nil {
		return &ResponseHandler{err: fmt.Errorf("create request: %w", err)}
	}

	// 合并 Headers（Client 全局 + 请求级别）
	for k, v := range rb.client.opts.Headers {
		req.Header[k] = v
	}
	for k, v := range rb.headers {
		req.Header[k] = v
	}

	// 执行请求
	resp, err := rb.client.httpClient.Do(req)
	return &ResponseHandler{resp: resp, err: err}
}

// 便捷方法
func (c *Client) PostJSON(ctx context.Context, url string, reqBody, respBody any) error {
	return c.NewRequest().
		Post(url).
		JSONBody(reqBody).
		Do(ctx).
		Into(respBody)
}

func (c *Client) GetJSON(ctx context.Context, url string, respBody any) error {
	return c.NewRequest().
		Get(url).
		Do(ctx).
		Into(respBody)
}

func (c *Client) PutJSON(ctx context.Context, url string, reqBody, respBody any) error {
	return c.NewRequest().
		Put(url).
		JSONBody(reqBody).
		Do(ctx).
		Into(respBody)
}

func (c *Client) DeleteJSON(ctx context.Context, url string, respBody any) error {
	return c.NewRequest().
		Delete(url).
		Do(ctx).
		Into(respBody)
}
