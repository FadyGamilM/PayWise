// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: service_paywise.proto

package pb

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

// PaywiseClient is the client API for Paywise service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaywiseClient interface {
	SignupUser(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error)
	LoginUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type paywiseClient struct {
	cc grpc.ClientConnInterface
}

func NewPaywiseClient(cc grpc.ClientConnInterface) PaywiseClient {
	return &paywiseClient{cc}
}

func (c *paywiseClient) SignupUser(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error) {
	out := new(SignupResponse)
	err := c.cc.Invoke(ctx, "/pb.Paywise/SignupUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paywiseClient) LoginUser(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/pb.Paywise/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaywiseServer is the server API for Paywise service.
// All implementations must embed UnimplementedPaywiseServer
// for forward compatibility
type PaywiseServer interface {
	SignupUser(context.Context, *SignupRequest) (*SignupResponse, error)
	LoginUser(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedPaywiseServer()
}

// UnimplementedPaywiseServer must be embedded to have forward compatible implementations.
type UnimplementedPaywiseServer struct {
}

func (UnimplementedPaywiseServer) SignupUser(context.Context, *SignupRequest) (*SignupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignupUser not implemented")
}
func (UnimplementedPaywiseServer) LoginUser(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedPaywiseServer) mustEmbedUnimplementedPaywiseServer() {}

// UnsafePaywiseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaywiseServer will
// result in compilation errors.
type UnsafePaywiseServer interface {
	mustEmbedUnimplementedPaywiseServer()
}

func RegisterPaywiseServer(s grpc.ServiceRegistrar, srv PaywiseServer) {
	s.RegisterService(&Paywise_ServiceDesc, srv)
}

func _Paywise_SignupUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaywiseServer).SignupUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Paywise/SignupUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaywiseServer).SignupUser(ctx, req.(*SignupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Paywise_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaywiseServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Paywise/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaywiseServer).LoginUser(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Paywise_ServiceDesc is the grpc.ServiceDesc for Paywise service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Paywise_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Paywise",
	HandlerType: (*PaywiseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignupUser",
			Handler:    _Paywise_SignupUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _Paywise_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_paywise.proto",
}
