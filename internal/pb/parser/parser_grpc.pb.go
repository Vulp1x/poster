// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package parser

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ParserClient is the client API for Parser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParserClient interface {
	// сохраняем ботов
	SaveBots(ctx context.Context, in *SaveBotsRequest, opts ...grpc.CallOption) (*SaveBotsResponse, error)
}

type parserClient struct {
	cc grpc.ClientConnInterface
}

func NewParserClient(cc grpc.ClientConnInterface) ParserClient {
	return &parserClient{cc}
}

func (c *parserClient) SaveBots(ctx context.Context, in *SaveBotsRequest, opts ...grpc.CallOption) (*SaveBotsResponse, error) {
	out := new(SaveBotsResponse)
	err := c.cc.Invoke(ctx, "/parser.Parser/SaveBots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParserServer is the server API for Parser service.
// All implementations must embed UnimplementedParserServer
// for forward compatibility
type ParserServer interface {
	// сохраняем ботов
	SaveBots(context.Context, *SaveBotsRequest) (*SaveBotsResponse, error)
	mustEmbedUnimplementedParserServer()
}

// UnimplementedParserServer must be embedded to have forward compatible implementations.
type UnimplementedParserServer struct {
}

func (UnimplementedParserServer) SaveBots(context.Context, *SaveBotsRequest) (*SaveBotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveBots not implemented")
}
func (UnimplementedParserServer) mustEmbedUnimplementedParserServer() {}

// UnsafeParserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParserServer will
// result in compilation errors.
type UnsafeParserServer interface {
	mustEmbedUnimplementedParserServer()
}

func RegisterParserServer(s grpc.ServiceRegistrar, srv ParserServer) {
	s.RegisterService(&Parser_ServiceDesc, srv)
}

func _Parser_SaveBots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveBotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParserServer).SaveBots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/parser.Parser/SaveBots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParserServer).SaveBots(ctx, req.(*SaveBotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Parser_ServiceDesc is the grpc.ServiceDesc for Parser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Parser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "parser.Parser",
	HandlerType: (*ParserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveBots",
			Handler:    _Parser_SaveBots_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "parser/parser.proto",
}
