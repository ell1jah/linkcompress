// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// LinkCompresserClient is the client API for LinkCompresser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinkCompresserClient interface {
	Compress(ctx context.Context, in *Link, opts ...grpc.CallOption) (*Link, error)
	Original(ctx context.Context, in *Link, opts ...grpc.CallOption) (*Link, error)
}

type linkCompresserClient struct {
	cc grpc.ClientConnInterface
}

func NewLinkCompresserClient(cc grpc.ClientConnInterface) LinkCompresserClient {
	return &linkCompresserClient{cc}
}

func (c *linkCompresserClient) Compress(ctx context.Context, in *Link, opts ...grpc.CallOption) (*Link, error) {
	out := new(Link)
	err := c.cc.Invoke(ctx, "/proto.LinkCompresser/Compress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkCompresserClient) Original(ctx context.Context, in *Link, opts ...grpc.CallOption) (*Link, error) {
	out := new(Link)
	err := c.cc.Invoke(ctx, "/proto.LinkCompresser/Original", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkCompresserServer is the server API for LinkCompresser service.
// All implementations must embed UnimplementedLinkCompresserServer
// for forward compatibility
type LinkCompresserServer interface {
	Compress(context.Context, *Link) (*Link, error)
	Original(context.Context, *Link) (*Link, error)
	mustEmbedUnimplementedLinkCompresserServer()
}

// UnimplementedLinkCompresserServer must be embedded to have forward compatible implementations.
type UnimplementedLinkCompresserServer struct {
}

func (UnimplementedLinkCompresserServer) Compress(context.Context, *Link) (*Link, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Compress not implemented")
}
func (UnimplementedLinkCompresserServer) Original(context.Context, *Link) (*Link, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Original not implemented")
}
func (UnimplementedLinkCompresserServer) mustEmbedUnimplementedLinkCompresserServer() {}

// UnsafeLinkCompresserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LinkCompresserServer will
// result in compilation errors.
type UnsafeLinkCompresserServer interface {
	mustEmbedUnimplementedLinkCompresserServer()
}

func RegisterLinkCompresserServer(s grpc.ServiceRegistrar, srv LinkCompresserServer) {
	s.RegisterService(&LinkCompresser_ServiceDesc, srv)
}

func _LinkCompresser_Compress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Link)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkCompresserServer).Compress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LinkCompresser/Compress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkCompresserServer).Compress(ctx, req.(*Link))
	}
	return interceptor(ctx, in, info, handler)
}

func _LinkCompresser_Original_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Link)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkCompresserServer).Original(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LinkCompresser/Original",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkCompresserServer).Original(ctx, req.(*Link))
	}
	return interceptor(ctx, in, info, handler)
}

// LinkCompresser_ServiceDesc is the grpc.ServiceDesc for LinkCompresser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LinkCompresser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LinkCompresser",
	HandlerType: (*LinkCompresserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Compress",
			Handler:    _LinkCompresser_Compress_Handler,
		},
		{
			MethodName: "Original",
			Handler:    _LinkCompresser_Original_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "link.proto",
}