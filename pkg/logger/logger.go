package logger

import (
	"context"
	"io"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global logger instance.
	global       *otelzap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
)

func init() {
	SetLogger(otelzap.New(
		New(defaultLevel),
		otelzap.WithTraceIDField(true),
		otelzap.WithMinLevel(zapcore.DebugLevel),
	).Sugar())
}

// New creates new *zap.SugaredLogger with standard EncoderConfig
// if lvl == nil, global AtomicLevel will be used
func New(level zapcore.LevelEnabler, options ...zap.Option) *zap.Logger {
	return NewWithSink(level, os.Stdout, options...)
}

func NewWithSink(level zapcore.LevelEnabler, sink io.Writer, options ...zap.Option) *zap.Logger {
	if level == nil {
		level = defaultLevel
	}

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(sink),
			level,
		),
		options...,
	)
}

// Level returns current global logger level
func Level() zapcore.Level {
	return defaultLevel.Level()
}

// SetLevel sets level for global logger
func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

// Logger returns current global logger.
func Logger() *otelzap.SugaredLogger {
	return global
}

// SetLogger sets global used logger. This function is not thread-safe.
func SetLogger(l *otelzap.SugaredLogger) {
	global = l
	otelzap.ReplaceGlobals(l.Desugar())
}

// Below listed all logging functions
// Suffix meaning:
// * No suffix, e.g. DebugLevel()   - log concatenated args
// * f,         e.g. Debugf()  - log using format string
// * KV,        e.g. DebugKV() - log key-values, odd args are keys, even – values
//

func DebugKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	FromContext(ctx).DebugwContext(ctx, msg, keysAndValues...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).DebugfContext(ctx, format, args...)
}

func InfoKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	FromContext(ctx).InfowContext(ctx, msg, keysAndValues...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).InfofContext(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).WarnfContext(ctx, format, args...)
}

func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).WarnwContext(ctx, message, kvs...)
}

func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).ErrorwContext(ctx, message, kvs...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).ErrorfContext(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).FatalfContext(ctx, format, args...)
}

func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Fatalw(message, kvs...)
}

func Panic(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		ctx = WithKV(ctx, "trace_id", traceID.String())
	}

	FromContext(ctx).Panicf(format, args...)
}

func PanicKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Panicw(message, kvs...)
}
