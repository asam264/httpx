package httpx

import (
	"net"
	"net/http"
	"time"
)

// defaultTransport 生产级 Transport 配置
func defaultTransport() *http.Transport {
	return &http.Transport{
		// 连接池配置
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 10,               // 每个 host 最大空闲连接
		MaxConnsPerHost:     0,                // 每个 host 最大连接数（0=无限制）
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时

		// TCP 配置
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时
			KeepAlive: 30 * time.Second, // TCP KeepAlive
		}).DialContext,

		// TLS 握手超时
		TLSHandshakeTimeout: 10 * time.Second,

		// 响应头超时
		ResponseHeaderTimeout: 10 * time.Second,

		// 期望 100-continue 超时
		ExpectContinueTimeout: 1 * time.Second,

		// 强制使用 HTTP/2
		ForceAttemptHTTP2: true,

		// 禁用压缩（如果需要手动控制）
		// DisableCompression: false,
	}
}

// CustomTransport 自定义 Transport 构建器
type TransportBuilder struct {
	transport *http.Transport
}

func NewTransport() *TransportBuilder {
	return &TransportBuilder{
		transport: defaultTransport(),
	}
}

func (tb *TransportBuilder) MaxIdleConns(n int) *TransportBuilder {
	tb.transport.MaxIdleConns = n
	return tb
}

func (tb *TransportBuilder) MaxIdleConnsPerHost(n int) *TransportBuilder {
	tb.transport.MaxIdleConnsPerHost = n
	return tb
}

func (tb *TransportBuilder) IdleConnTimeout(d time.Duration) *TransportBuilder {
	tb.transport.IdleConnTimeout = d
	return tb
}

func (tb *TransportBuilder) DialTimeout(d time.Duration) *TransportBuilder {
	tb.transport.DialContext = (&net.Dialer{
		Timeout:   d,
		KeepAlive: 30 * time.Second,
	}).DialContext
	return tb
}

func (tb *TransportBuilder) Build() *http.Transport {
	return tb.transport
}
