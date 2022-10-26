package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/pkg/logger"
)

const requestIDHeaderKey = "X-Request-ID"

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
	ctx := req.Context()
	if req.Header.Get(requestIDHeaderKey) == "" {
		req.Header.Set(requestIDHeaderKey, uuid.NewString())
	}

	logger.Infof(ctx, "sending request to %s, request_id: '%s'", req.URL, req.Header.Get(requestIDHeaderKey))
	startedAt := time.Now()

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	var statusCode = -666
	if res != nil {
		statusCode = res.StatusCode
	}

	logger.Infof(ctx, "got response in %s, status code: %d", time.Since(startedAt), statusCode)

	return
}
