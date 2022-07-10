// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: user/proto/account.proto

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

// AccountManagementClient is the client API for AccountManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountManagementClient interface {
	AddAccount(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountMessage, error)
	GetAccountsByFilter(ctx context.Context, in *GetAccountsByFilterRequest, opts ...grpc.CallOption) (*AccountsResponse, error)
}

type accountManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountManagementClient(cc grpc.ClientConnInterface) AccountManagementClient {
	return &accountManagementClient{cc}
}

func (c *accountManagementClient) AddAccount(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountMessage, error) {
	out := new(AccountMessage)
	err := c.cc.Invoke(ctx, "/product.AccountManagement/AddAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountManagementClient) GetAccountsByFilter(ctx context.Context, in *GetAccountsByFilterRequest, opts ...grpc.CallOption) (*AccountsResponse, error) {
	out := new(AccountsResponse)
	err := c.cc.Invoke(ctx, "/product.AccountManagement/GetAccountsByFilter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountManagementServer is the server API for AccountManagement service.
// All implementations must embed UnimplementedAccountManagementServer
// for forward compatibility
type AccountManagementServer interface {
	AddAccount(context.Context, *AccountRequest) (*AccountMessage, error)
	GetAccountsByFilter(context.Context, *GetAccountsByFilterRequest) (*AccountsResponse, error)
	mustEmbedUnimplementedAccountManagementServer()
}

// UnimplementedAccountManagementServer must be embedded to have forward compatible implementations.
type UnimplementedAccountManagementServer struct {
}

func (UnimplementedAccountManagementServer) AddAccount(context.Context, *AccountRequest) (*AccountMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAccount not implemented")
}
func (UnimplementedAccountManagementServer) GetAccountsByFilter(context.Context, *GetAccountsByFilterRequest) (*AccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountsByFilter not implemented")
}
func (UnimplementedAccountManagementServer) mustEmbedUnimplementedAccountManagementServer() {}

// UnsafeAccountManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountManagementServer will
// result in compilation errors.
type UnsafeAccountManagementServer interface {
	mustEmbedUnimplementedAccountManagementServer()
}

func RegisterAccountManagementServer(s grpc.ServiceRegistrar, srv AccountManagementServer) {
	s.RegisterService(&AccountManagement_ServiceDesc, srv)
}

func _AccountManagement_AddAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountManagementServer).AddAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.AccountManagement/AddAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountManagementServer).AddAccount(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountManagement_GetAccountsByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountsByFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountManagementServer).GetAccountsByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.AccountManagement/GetAccountsByFilter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountManagementServer).GetAccountsByFilter(ctx, req.(*GetAccountsByFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountManagement_ServiceDesc is the grpc.ServiceDesc for AccountManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "product.AccountManagement",
	HandlerType: (*AccountManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddAccount",
			Handler:    _AccountManagement_AddAccount_Handler,
		},
		{
			MethodName: "GetAccountsByFilter",
			Handler:    _AccountManagement_GetAccountsByFilter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/proto/account.proto",
}
