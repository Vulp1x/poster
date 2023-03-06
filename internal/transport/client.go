package transport

import (
	"net/http"
	"time"
)

func InitHTTPClient() *http.Client {

	return &http.Client{
		Transport: &loggingRoundTripper{Proxied: http.DefaultTransport},
		Timeout:   200 * time.Second,
	}
}

func ProxyingHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = FromContext()

	return &http.Client{
		Transport: &loggingRoundTripper{Proxied: transport},
		Timeout:   timeout,
	}
}
