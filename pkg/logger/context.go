package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota
)

// ToContext returns new context with specified sugared logger inside.
func ToContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

// FromContext returns logger from context if set. Otherwise returns global `global` logger.
// In both cases returned logger is populated with `trace_id` & `span_id`.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	log := global

	if logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger); ok {
		log = logger
	}

	return log
}

// WithName creates a child logger from provided context and names it.
// Resulting log entries will be enriched with name property.
// Childs of child logger will stake names.
//
// Example:
//
//	ctx := WithName(ctx, "GetApples") -> "GetApples"
//	ctx = WithName(ctx, "AppleManager") -> "GetApples.AppleManager"
//	ctx = WithName(ctx, "DB") -> "GetApples.AppleManager.DB"
func WithName(ctx context.Context, name string) context.Context {
	return ToContext(ctx, FromContext(ctx).Named(name))
}

// WithKV creates a child logger from provided context and sets custom metadata.
// It accepts key-value pairs, which will be passed to all child loggers.
//
// Example:
//
//	ctx := WithKV(ctx,
//		"filter", rqFilter,
//	)
func WithKV(ctx context.Context, key string, value interface{}) context.Context {
	log := FromContext(ctx).With(key, value)
	return ToContext(ctx, log)
}

// WithFields creates a child logger from provided and sets typed fields metadata.
//
// Example:
//
//	ctx := WithFields(ctx,
//		logger.Fields{
//			"kafka-topic": topic,
//			"kafka-partition": partition,
//		})
func WithFields(ctx context.Context, fields Fields) context.Context {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	log := FromContext(ctx).
		Desugar().
		With(zapFields...).
		Sugar()
	return ToContext(ctx, log)
}
