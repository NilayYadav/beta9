// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.1
// source: shell.proto

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

// ShellServiceClient is the client API for ShellService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShellServiceClient interface {
	StartShell(ctx context.Context, opts ...grpc.CallOption) (ShellService_StartShellClient, error)
}

type shellServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShellServiceClient(cc grpc.ClientConnInterface) ShellServiceClient {
	return &shellServiceClient{cc}
}

func (c *shellServiceClient) StartShell(ctx context.Context, opts ...grpc.CallOption) (ShellService_StartShellClient, error) {
	stream, err := c.cc.NewStream(ctx, &ShellService_ServiceDesc.Streams[0], "/shell.ShellService/StartShell", opts...)
	if err != nil {
		return nil, err
	}
	x := &shellServiceStartShellClient{stream}
	return x, nil
}

type ShellService_StartShellClient interface {
	Send(*ShellRequest) error
	Recv() (*ShellResponse, error)
	grpc.ClientStream
}

type shellServiceStartShellClient struct {
	grpc.ClientStream
}

func (x *shellServiceStartShellClient) Send(m *ShellRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *shellServiceStartShellClient) Recv() (*ShellResponse, error) {
	m := new(ShellResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShellServiceServer is the server API for ShellService service.
// All implementations must embed UnimplementedShellServiceServer
// for forward compatibility
type ShellServiceServer interface {
	StartShell(ShellService_StartShellServer) error
	mustEmbedUnimplementedShellServiceServer()
}

// UnimplementedShellServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShellServiceServer struct {
}

func (UnimplementedShellServiceServer) StartShell(ShellService_StartShellServer) error {
	return status.Errorf(codes.Unimplemented, "method StartShell not implemented")
}
func (UnimplementedShellServiceServer) mustEmbedUnimplementedShellServiceServer() {}

// UnsafeShellServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShellServiceServer will
// result in compilation errors.
type UnsafeShellServiceServer interface {
	mustEmbedUnimplementedShellServiceServer()
}

func RegisterShellServiceServer(s grpc.ServiceRegistrar, srv ShellServiceServer) {
	s.RegisterService(&ShellService_ServiceDesc, srv)
}

func _ShellService_StartShell_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ShellServiceServer).StartShell(&shellServiceStartShellServer{stream})
}

type ShellService_StartShellServer interface {
	Send(*ShellResponse) error
	Recv() (*ShellRequest, error)
	grpc.ServerStream
}

type shellServiceStartShellServer struct {
	grpc.ServerStream
}

func (x *shellServiceStartShellServer) Send(m *ShellResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *shellServiceStartShellServer) Recv() (*ShellRequest, error) {
	m := new(ShellRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShellService_ServiceDesc is the grpc.ServiceDesc for ShellService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShellService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shell.ShellService",
	HandlerType: (*ShellServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartShell",
			Handler:       _ShellService_StartShell_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "shell.proto",
}