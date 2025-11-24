package httpx

import (
	"net/http"
	"time"
)

type Options struct {
	Timeout         time.Duration
	MaxRetries      int
	RetryMinBackoff time.Duration
	RetryMaxBackoff time.Duration
	RetryIf         RetryConditionFunc
	Headers         http.Header
	Middlewares     []Middleware
	Transport       *http.Transport
}

type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Timeout:         10 * time.Second,
		MaxRetries:      0, // 默认不重试
		RetryMinBackoff: 100 * time.Millisecond,
		RetryMaxBackoff: 5 * time.Second,
		Headers:         make(http.Header),
		Middlewares:     []Middleware{},
	}
}

func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

func WithRetry(maxRetries int) Option {
	return func(o *Options) {
		o.MaxRetries = maxRetries
	}
}

func WithRetryBackoff(min, max time.Duration) Option {
	return func(o *Options) {
		o.RetryMinBackoff = min
		o.RetryMaxBackoff = max
	}
}

func WithRetryIf(fn RetryConditionFunc) Option {
	return func(o *Options) {
		o.RetryIf = fn
	}
}

func WithHeader(key, value string) Option {
	return func(o *Options) {
		o.Headers.Set(key, value)
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(o *Options) {
		for k, v := range headers {
			o.Headers.Set(k, v)
		}
	}
}

func WithMiddleware(mw ...Middleware) Option {
	return func(o *Options) {
		o.Middlewares = append(o.Middlewares, mw...)
	}
}

func WithTransport(t *http.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}
