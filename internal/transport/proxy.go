package transport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ctxKey string

var proxyKey = ctxKey("proxy_url")

// ContextWithProxy adds transport to context
func ContextWithProxy(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, proxyKey, u)
}

// FromContext use transport from request context
func FromContext() func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		ctx := req.Context()

		proxy, ok := ctx.Value(proxyKey).(*url.URL)
		if !ok {
			return nil, fmt.Errorf("failed to get transport from ctx: %+v", ctx)
		}

		return proxy, nil
	}
}

func roundRobin() func(*http.Request) (*url.URL, error) {
	var i = 0
	return func(request *http.Request) (*url.URL, error) {
		i++
		if i%2 == 0 {
			return &url.URL{
				Scheme: "http",
				User:   url.UserPassword("dmitrijkholodkov7815", "21e49b"),
				Host:   "109.248.7.160:10534",
			}, nil
		}

		return &url.URL{
			Scheme: "http",
			User:   url.UserPassword("dmitrijkholodkov7815", "21e49b"),
			Host:   "193.23.50.220:10117",
		}, nil
	}
}
