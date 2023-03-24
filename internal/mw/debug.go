package mw

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/inst-api/poster/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	goahttp "goa.design/goa/v3/http"
)

// responseDupper tees the response to a buffer and a response writer.
type responseDupper struct {
	http.ResponseWriter
	Buffer      *bytes.Buffer
	wroteHeader bool
	code        int
	bytes       int
	startedAt   time.Time
}

// RequestLoggerWithDebug returns a debug middleware which prints detailed information about
// incoming requests and outgoing responses including all headers, parameters
// and bodies.
func RequestLoggerWithDebug(mux goahttp.Muxer, debug bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			f := StructuredLogger{debug: debug}
			ctx, entry := f.NewLogEntry(r)
			r = r.WithContext(ctx)

			dupper := &responseDupper{ResponseWriter: rw, Buffer: &bytes.Buffer{}, startedAt: time.Now()}

			t1 := time.Now()
			defer func() {
				entry.Write(dupper.Status(), dupper.BytesWritten(), dupper.Header(), time.Since(t1), dupper)
			}()

			rawTraceID := trace.SpanContextFromContext(r.Context()).TraceID()
			if rawTraceID.IsValid() {
				dupper.ResponseWriter.Header().Add("X-Trace-ID", rawTraceID.String())
			} else {
				logger.Infof(r.Context(), "failed to get trace_id from '%v' (type %T)", rawTraceID, rawTraceID)
			}

			h.ServeHTTP(dupper, chimw.WithLogEntry(r, entry))
		})
	}
}

// Hijack supports the http.Hijacker interface.
func (r *responseDupper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("debug middleware: inner ResponseWriter cannot be hijacked: %T", r.ResponseWriter)
}

func (r *responseDupper) WriteHeader(code int) {
	if !r.wroteHeader {
		r.code = code
		r.wroteHeader = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *responseDupper) Write(buf []byte) (int, error) {
	r.maybeWriteHeader()
	r.Buffer.Write(buf)
	n, err := r.ResponseWriter.Write(buf)
	r.bytes += n
	return n, err
}

func (r *responseDupper) maybeWriteHeader() {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
}

func (r *responseDupper) Status() int {
	return r.code
}

func (r *responseDupper) BytesWritten() int {
	return r.bytes
}
