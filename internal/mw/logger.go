package mw

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/inst-api/poster/internal/tracer"
	"github.com/inst-api/poster/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
)

var _ middleware.LogEntry = &StructuredLoggerEntry{}

// StructuredLogger implements logger.
type StructuredLogger struct {
	debug bool
}

// NewLogEntry creates logger.
func (l *StructuredLogger) NewLogEntry(r *http.Request) (context.Context, middleware.LogEntry) {
	ctx, span := tracer.Start(r.Context(), "body_reader")
	defer span.End()

	var entry middleware.LogEntry
	if l.debug {
		entry = &DebugStructuredLoggerEntry{ctx: ctx}
	} else {
		entry = &StructuredLoggerEntry{ctx: ctx}
	}

	if l.debug {
		span.SetAttributes(attribute.String("content_length", byteCount(r.ContentLength)))
		if r.ContentLength < 10_000 {
			// Request body
			b, err := io.ReadAll(r.Body)
			if err != nil {
				b = []byte("failed to read body: " + err.Error())
			}

			if err = r.Body.Close(); err != nil {
				logger.Errorf(ctx, "failed to close body: %v", err)
			}

			span.SetAttributes(attribute.String("body", string(b)))
			r.Body = io.NopCloser(bytes.NewBuffer(b))
		}

	}

	logger.DebugKV(ctx, "request started")

	return ctx, entry
}

// StructuredLoggerEntry ...
type StructuredLoggerEntry struct {
	ctx context.Context
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.ctx = logger.WithFields(l.ctx, logger.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) * float64(time.Nanosecond) / float64(time.Millisecond),
	})

	logger.Infof(l.ctx, "request complete")
}

// Panic implemetents logEntry method.
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.ctx = logger.WithFields(l.ctx, logger.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// DebugStructuredLoggerEntry ...
type DebugStructuredLoggerEntry struct {
	ctx context.Context
}

func (l *DebugStructuredLoggerEntry) Write(status, bytesWriten int, header http.Header, elapsed time.Duration, extra interface{}) {
	ctx, span := tracer.Start(l.ctx, "response_saver")
	defer span.End()

	fields := logger.Fields{
		"resp.status": status, "resp.bytes_length": byteCount(int64(bytesWriten)),
		"resp.elapsed_ms": float64(elapsed.Nanoseconds()) * float64(time.Nanosecond) / float64(time.Millisecond),
	}

	ww, ok := extra.(*responseDupper)
	if ok {
		fields["resp.elapsed_ms"] = float64(time.Since(ww.startedAt).Nanoseconds()) * float64(time.Nanosecond) / float64(time.Millisecond)
		fields["resp.headers"] = header
	}

	if ww.Buffer.Len() < 10000 {
		fields["resp.body"] = ww.Buffer.String()
	}

	logger.DebugKV(logger.WithFields(ctx, fields), "request completed")
}

// Panic implemetents logEntry method.
func (l *DebugStructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.ctx = logger.WithFields(l.ctx, logger.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// InternalError helper for sending error.
func InternalError(
	w http.ResponseWriter,
	r *http.Request,
	msg string, args ...interface{}) {
	logger.Errorf(r.Context(), msg, args...)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// BadRequest helper for sending error.
func BadRequest(
	w http.ResponseWriter,
	r *http.Request,
	msg string, args ...interface{}) {
	logger.Errorf(r.Context(), msg, args...)

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// Error helper for sending error.
func Error(
	w http.ResponseWriter,
	r *http.Request,
	code int,
	msg string, args ...interface{}) {
	logger.Errorf(r.Context(), msg, args...)

	http.Error(w, http.StatusText(code), code)
}

func byteCount(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
