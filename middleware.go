package httpx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.RoundTripper) http.RoundTripper
type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// middlewareTransport 中间件传输层
type middlewareTransport struct {
	base http.RoundTripper
}

func (mt *middlewareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return mt.base.RoundTrip(req)
}

// buildMiddlewareChain 构建中间件链
func buildMiddlewareChain(base http.RoundTripper, opts *Options) http.RoundTripper {
	// 从后往前包装
	rt := base

	// 1. 重试层（最内层）
	if opts.MaxRetries > 0 {
		retryIf := opts.RetryIf
		if retryIf == nil {
			retryIf = DefaultRetryIf
		}
		rt = &retryTransport{
			base:       rt,
			maxRetries: opts.MaxRetries,
			minBackoff: opts.RetryMinBackoff,
			maxBackoff: opts.RetryMaxBackoff,
			retryIf:    retryIf,
		}
	}

	// 2. 用户自定义中间件
	for i := len(opts.Middlewares) - 1; i >= 0; i-- {
		rt = opts.Middlewares[i](rt)
	}

	return &middlewareTransport{base: rt}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()

			log.Printf("[HTTP] --> %s %s", req.Method, req.URL.String())

			resp, err := next.RoundTrip(req)

			duration := time.Since(start)
			if err != nil {
				log.Printf("[HTTP] <-- %s %s | ERROR: %v | %dms",
					req.Method, req.URL.String(), err, duration.Milliseconds())
			} else {
				log.Printf("[HTTP] <-- %s %s | %d | %dms",
					req.Method, req.URL.String(), resp.StatusCode, duration.Milliseconds())
			}

			return resp, err
		})
	}
}

// MetricsMiddleware Prometheus 指标中间件（示例）
func MetricsMiddleware(serviceName string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := next.RoundTrip(req)
			duration := time.Since(start)

			// 这里可以集成 Prometheus
			status := "error"
			if resp != nil {
				status = fmt.Sprintf("%d", resp.StatusCode)
			}

			// prometheus.RecordHTTPRequest(serviceName, req.Method, status, duration)
			_ = status
			_ = duration

			return resp, err
		})
	}
}

// TimeoutMiddleware 请求级超时中间件
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return next.RoundTrip(req.WithContext(ctx))
		})
	}
}
