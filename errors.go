package httpx

import (
	"errors"
	"fmt"
)

// HTTPError HTTP 错误响应
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e *HTTPError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("http %d: %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Status)
}

// IsHTTPError 判断是否为 HTTP 错误
func IsHTTPError(err error) bool {
	var httpErr *HTTPError
	return errors.As(err, &httpErr)
}

// GetHTTPError 提取 HTTP 错误
func GetHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	ok := errors.As(err, &httpErr)
	return httpErr, ok
}

// IsStatusCode 判断错误是否为特定状态码
func IsStatusCode(err error, code int) bool {
	httpErr, ok := GetHTTPError(err)
	return ok && httpErr.StatusCode == code
}

// IsTimeout 判断是否为超时错误
func IsTimeout(err error) bool {
	type timeout interface {
		Timeout() bool
	}
	var t timeout
	return errors.As(err, &t) && t.Timeout()
}
