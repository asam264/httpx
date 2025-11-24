package httpx

import (
	"context"
	"net/http"
	"sync"
	"time"
)

var (
	defaultClient     *Client
	defaultClientOnce sync.Once
)

type Client struct {
	httpClient *http.Client
	opts       *Options
	baseURL    string
}

// New 创建新客户端
func New(options ...Option) *Client {
	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}

	transport := opts.Transport
	if transport == nil {
		transport = defaultTransport()
	}

	httpClient := &http.Client{
		Timeout:   opts.Timeout,
		Transport: buildMiddlewareChain(transport, opts),
	}

	return &Client{
		httpClient: httpClient,
		opts:       opts,
	}
}

// SetDefault 设置全局默认客户端
func SetDefault(client *Client) {
	defaultClient = client
}

func getDefaultClient() *Client {
	defaultClientOnce.Do(func() {
		defaultClient = New()
	})
	return defaultClient
}

// 全局快捷方法
func PostJSON(ctx context.Context, url string, reqBody, respBody any) error {
	return getDefaultClient().PostJSON(ctx, url, reqBody, respBody)
}

func GetJSON(ctx context.Context, url string, respBody any) error {
	return getDefaultClient().GetJSON(ctx, url, respBody)
}

// WithBaseURL 设置基础 URL
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

// WithTimeout 修改超时（返回新实例）
func (c *Client) WithTimeout(d time.Duration) *Client {
	newOpts := *c.opts
	newOpts.Timeout = d
	return New(func(o *Options) { *o = newOpts })
}

// WithRetry 配置重试
func (c *Client) WithRetry(maxRetries int) *Client {
	c.opts.MaxRetries = maxRetries
	return c
}

func (c *Client) WithRetryBackoff(min, max time.Duration) *Client {
	c.opts.RetryMinBackoff = min
	c.opts.RetryMaxBackoff = max
	return c
}

func (c *Client) WithRetryIf(fn RetryConditionFunc) *Client {
	c.opts.RetryIf = fn
	return c
}

func (c *Client) WithHeader(key, value string) *Client {
	c.opts.Headers.Set(key, value)
	return c
}

func (c *Client) WithMiddleware(mw ...Middleware) *Client {
	c.opts.Middlewares = append(c.opts.Middlewares, mw...)
	// 重建 Transport 链
	c.httpClient.Transport = buildMiddlewareChain(
		c.httpClient.Transport.(*middlewareTransport).base,
		c.opts,
	)
	return c
}

// NewRequest 创建请求构建器
func (c *Client) NewRequest() *RequestBuilder {
	return &RequestBuilder{
		client:  c,
		headers: make(http.Header),
		queries: make(map[string][]string),
	}
}
