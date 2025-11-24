package httpx

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

type RetryConditionFunc func(resp *http.Response, err error) bool

// DefaultRetryIf 默认重试条件：网络错误 + 5xx + 429
func DefaultRetryIf(resp *http.Response, err error) bool {
	if err != nil {
		return true // 网络错误重试
	}
	if resp == nil {
		return false
	}
	// 重试服务端错误和限流
	return resp.StatusCode >= 500 || resp.StatusCode == 429
}

// retryTransport 重试传输层
type retryTransport struct {
	base       http.RoundTripper
	maxRetries int
	minBackoff time.Duration
	maxBackoff time.Duration
	retryIf    RetryConditionFunc
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.maxRetries; attempt++ {
		// 检查 context 是否已取消
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// 执行请求
		resp, err = rt.base.RoundTrip(req)

		// 判断是否需要重试
		shouldRetry := rt.retryIf(resp, err)
		if !shouldRetry || attempt == rt.maxRetries {
			return resp, err
		}

		// 计算退避时间
		backoff := rt.calculateBackoff(attempt)

		// 等待退避时间
		select {
		case <-time.After(backoff):
			// 继续重试
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		// 关闭之前的响应体（如果有）
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	return resp, err
}

// calculateBackoff 计算退避时间（指数退避 + jitter）
func (rt *retryTransport) calculateBackoff(attempt int) time.Duration {
	// 指数退避: min * 2^attempt
	backoff := float64(rt.minBackoff) * math.Pow(2, float64(attempt))

	// 限制最大值
	if backoff > float64(rt.maxBackoff) {
		backoff = float64(rt.maxBackoff)
	}

	// 添加 jitter（±25%）
	jitter := backoff * 0.25 * (rand.Float64()*2 - 1)
	backoff += jitter

	return time.Duration(backoff)
}
