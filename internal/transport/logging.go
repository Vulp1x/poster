package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var elapsedTimeKey = ctxKey("elapsed_time")

// contextWithElapsedTime adds request elapsed time to context
func contextWithElapsedTime(ctx context.Context, duration time.Duration) context.Context {
	return context.WithValue(ctx, elapsedTimeKey, duration)
}

// This type implements the http.RoundTripper interface
type loggingRoundTripper struct {
	Proxied http.RoundTripper
}

// RoundTrip add started_at request field
func (lrt loggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// Do "before sending requests" actions here.
	fmt.Printf("Sending request to %v\n", req.URL)
	startedAt := time.Now()

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	req.WithContext(contextWithElapsedTime(req.Context(), time.Since(startedAt)))

	return
}
