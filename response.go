package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseHandler struct {
	resp *http.Response
	err  error
}

// Into 解析 JSON 响应
func (rh *ResponseHandler) Into(v any) error {
	if rh.err != nil {
		return rh.err
	}

	defer rh.resp.Body.Close()

	// 检查状态码
	if rh.resp.StatusCode < 200 || rh.resp.StatusCode >= 300 {
		body, _ := io.ReadAll(rh.resp.Body)
		return &HTTPError{
			StatusCode: rh.resp.StatusCode,
			Status:     rh.resp.Status,
			Body:       body,
		}
	}

	// 解码 JSON
	if err := json.NewDecoder(rh.resp.Body).Decode(v); err != nil {
		return fmt.Errorf("decode json response: %w", err)
	}

	return nil
}

// Raw 返回原始响应（需要手动关闭 Body）
func (rh *ResponseHandler) Raw() (*http.Response, error) {
	return rh.resp, rh.err
}

// Bytes 读取响应体为字节
func (rh *ResponseHandler) Bytes() ([]byte, error) {
	if rh.err != nil {
		return nil, rh.err
	}
	defer rh.resp.Body.Close()

	if rh.resp.StatusCode < 200 || rh.resp.StatusCode >= 300 {
		body, _ := io.ReadAll(rh.resp.Body)
		return nil, &HTTPError{
			StatusCode: rh.resp.StatusCode,
			Status:     rh.resp.Status,
			Body:       body,
		}
	}

	return io.ReadAll(rh.resp.Body)
}

// String 读取响应体为字符串
func (rh *ResponseHandler) String() (string, error) {
	b, err := rh.Bytes()
	return string(b), err
}

// StatusCode 获取状态码
func (rh *ResponseHandler) StatusCode() int {
	if rh.resp == nil {
		return 0
	}
	return rh.resp.StatusCode
}
