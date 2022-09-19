package transport

import (
	"net/http"
	"time"
)

func InitHTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: FromContext()}, Timeout: 5 * time.Second}
}
