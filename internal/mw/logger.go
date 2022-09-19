package mw

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/inst-api/poster/pkg/logger"
)

var _ middleware.LogEntry = &StructuredLoggerEntry{}

// StructuredLogger implements logger.
type StructuredLogger struct {
	debug bool
}

// NewLogEntry creates logger.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {

	logFields := logger.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logFields["http_scheme"] = scheme
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entryCtx := logger.WithFields(r.Context(), logFields)

	var entry middleware.LogEntry
	if l.debug {
		entry = &DebugStructuredLoggerEntry{ctx: entryCtx}
	} else {
		entry = &StructuredLoggerEntry{ctx: entryCtx}
	}

	logCtx := entryCtx
	if l.debug {
		logDebugFields := logger.Fields{}

		// Request Headers
		keys := make([]string, len(r.Header))
		i := 0
		for k := range r.Header {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		logDebugFields["headers"] = keys

		// Request body
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			b = []byte("failed to read body: " + err.Error())
		}

		if len(b) > 0 {
			logDebugFields["req_body"] = b
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

		logCtx = logger.WithFields(logCtx, logDebugFields)
	}

	logger.Infof(logCtx, "request started")

	return entry
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

	fields := logger.Fields{
		"resp_status": status, "resp_bytes_length": bytesWriten,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) * float64(time.Nanosecond) / float64(time.Millisecond),
	}

	logger.Info(logger.WithFields(l.ctx, fields), "request complete")

	ww, ok := extra.(*responseDupper)
	if ok {
		fields["resp_body"] = ww.Buffer.String()
		fields["resp_elapsed_ms"] = float64(time.Since(ww.startedAt).Nanoseconds()) * float64(time.Nanosecond) / float64(time.Millisecond)
		fields["resp_headers"] = header
	}

	l.ctx = logger.WithFields(l.ctx, fields)

	logger.Debugf(l.ctx, "request complete")
}

// Panic implemetents logEntry method.
func (l *DebugStructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.ctx = logger.WithFields(l.ctx, logger.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

// GetLogEntry helper method in middleware chain.
func GetLogEntry(r *http.Request) *StructuredLoggerEntry {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)

	return entry
}

// LogEntrySetField helper method in middleware chain.
func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.ctx = logger.WithFields(entry.ctx, logger.Fields{key: value})
	}
}

// LogEntrySetFields helper method in middleware chain.
func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.ctx = logger.WithFields(entry.ctx, fields)
	}
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
