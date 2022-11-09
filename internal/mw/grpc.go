package mw

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/inst-api/poster/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// RequestIDMetadataKey is the key containing the request ID in the gRPC
	// metadata.
	RequestIDMetadataKey = "x-request-id"
)

func UnaryServerLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return UnaryLog(ctx, req, info, handler)
	}
}

func UnaryClientLog() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return unaryClientLog(ctx, method, req, reply, cc, invoker, opts...)
	}
}

// UnaryLog does the actual logging given the logger for unary methods.
func UnaryLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var reqID string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	reqID = MetadataValue(md, RequestIDMetadataKey)
	if reqID == "" {
		reqID = ShortID()
	}

	started := time.Now()

	ctx = logger.WithFields(ctx, logger.Fields{"req_id": reqID})

	startCtx := logger.WithFields(ctx, logger.Fields{"method": info.FullMethod, "request_length": byteCount(messageLength(req))})
	// before executing rpc
	logger.Info(startCtx, "request started")

	// invoke rpc
	resp, err = handler(ctx, req)

	// after executing rpc
	s, _ := status.FromError(err)

	afterCtx := logger.WithFields(ctx, logger.Fields{"status": s.Code(), "bytes": byteCount(messageLength(req)), "elapsed": time.Since(started).String()})
	logger.Info(afterCtx, "request completed")

	return resp, err
}

// UnaryClientLog does the actual logging given the logger for unary methods.
func unaryClientLog(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	reqID := ShortID()
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{RequestIDMetadataKey: reqID}))

	started := time.Now()

	ctx = logger.WithFields(ctx, logger.Fields{"req_id": reqID})

	startCtx := logger.WithFields(ctx, logger.Fields{"method": method, "request_length": byteCount(messageLength(req))})
	// before executing rpc
	logger.Info(startCtx, "request started")

	// invoke rpc
	err := invoker(ctx, method, req, reply, cc, opts...)

	// after executing rpc
	s, _ := status.FromError(err)

	afterCtx := logger.WithFields(ctx, logger.Fields{"status": s.Code(), "bytes": byteCount(messageLength(req)), "elapsed": time.Since(started).String()})
	logger.Info(afterCtx, "request completed")

	return err
}

// MetadataValue returns the first value for the given metadata key if
// key exists, else returns an empty string.
func MetadataValue(md metadata.MD, key string) string {
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func messageLength(msg interface{}) int64 {
	var length int64
	{
		if m, ok := msg.(proto.Message); ok {
			length = int64(proto.Size(m))
		}
	}
	return length
}

// UnaryRequestID returns a middleware for unary gRPC requests which
// initializes the request metadata with a unique value under the
// RequestIDMetadata key. Optionally, it uses the incoming "x-request-id"
// request metadata key, if present, with or without a length limit to use as
// request ID. The default behavior is to always generate a new ID.
//
// examples of use:
//
//	grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryRequestID()))
//
//	// enable options for using "x-request-id" metadata key with length limit.
//	grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryRequestID(
//	  middleware.UseXRequestIDMetadataOption(true),
//	  middleware.XRequestMetadataLimitOption(128))))
func UnaryRequestID() grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = generateRequestID(ctx)
		return handler(ctx, req)
	})
}

func generateRequestID(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	requestID := MetadataValue(md, RequestIDMetadataKey)
	if requestID == "" {
		requestID = ShortID()
	}

	md.Set(RequestIDMetadataKey, requestID)
	return metadata.NewIncomingContext(ctx, md)
}

// ShortID produces a " unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func ShortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
