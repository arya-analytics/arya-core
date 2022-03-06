// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chanchunkv1

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

// ChannelChunkServiceClient is the client API for ChannelChunkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelChunkServiceClient interface {
	CreateReplicas(ctx context.Context, opts ...grpc.CallOption) (ChannelChunkService_CreateReplicasClient, error)
	RetrieveReplicas(ctx context.Context, in *ChannelChunkServiceRetrieveReplicasRequest, opts ...grpc.CallOption) (ChannelChunkService_RetrieveReplicasClient, error)
	DeleteReplicas(ctx context.Context, in *ChannelChunkServiceDeleteReplicasRequest, opts ...grpc.CallOption) (*ChannelChunkServiceDeleteReplicasResponse, error)
}

type channelChunkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelChunkServiceClient(cc grpc.ClientConnInterface) ChannelChunkServiceClient {
	return &channelChunkServiceClient{cc}
}

func (c *channelChunkServiceClient) CreateReplicas(ctx context.Context, opts ...grpc.CallOption) (ChannelChunkService_CreateReplicasClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChannelChunkService_ServiceDesc.Streams[0], "/chanchunk.v1.ChannelChunkService/CreateReplicas", opts...)
	if err != nil {
		return nil, err
	}
	x := &channelChunkServiceCreateReplicasClient{stream}
	return x, nil
}

type ChannelChunkService_CreateReplicasClient interface {
	Send(*ChannelChunkServiceCreateReplicasRequest) error
	CloseAndRecv() (*ChannelChunkServiceCreateReplicasResponse, error)
	grpc.ClientStream
}

type channelChunkServiceCreateReplicasClient struct {
	grpc.ClientStream
}

func (x *channelChunkServiceCreateReplicasClient) Send(m *ChannelChunkServiceCreateReplicasRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *channelChunkServiceCreateReplicasClient) CloseAndRecv() (*ChannelChunkServiceCreateReplicasResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ChannelChunkServiceCreateReplicasResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *channelChunkServiceClient) RetrieveReplicas(ctx context.Context, in *ChannelChunkServiceRetrieveReplicasRequest, opts ...grpc.CallOption) (ChannelChunkService_RetrieveReplicasClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChannelChunkService_ServiceDesc.Streams[1], "/chanchunk.v1.ChannelChunkService/RetrieveReplicas", opts...)
	if err != nil {
		return nil, err
	}
	x := &channelChunkServiceRetrieveReplicasClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChannelChunkService_RetrieveReplicasClient interface {
	Recv() (*ChannelChunkServiceRetrieveReplicasResponse, error)
	grpc.ClientStream
}

type channelChunkServiceRetrieveReplicasClient struct {
	grpc.ClientStream
}

func (x *channelChunkServiceRetrieveReplicasClient) Recv() (*ChannelChunkServiceRetrieveReplicasResponse, error) {
	m := new(ChannelChunkServiceRetrieveReplicasResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *channelChunkServiceClient) DeleteReplicas(ctx context.Context, in *ChannelChunkServiceDeleteReplicasRequest, opts ...grpc.CallOption) (*ChannelChunkServiceDeleteReplicasResponse, error) {
	out := new(ChannelChunkServiceDeleteReplicasResponse)
	err := c.cc.Invoke(ctx, "/chanchunk.v1.ChannelChunkService/DeleteReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelChunkServiceServer is the server API for ChannelChunkService service.
// All implementations should embed UnimplementedChannelChunkServiceServer
// for forward compatibility
type ChannelChunkServiceServer interface {
	CreateReplicas(ChannelChunkService_CreateReplicasServer) error
	RetrieveReplicas(*ChannelChunkServiceRetrieveReplicasRequest, ChannelChunkService_RetrieveReplicasServer) error
	DeleteReplicas(context.Context, *ChannelChunkServiceDeleteReplicasRequest) (*ChannelChunkServiceDeleteReplicasResponse, error)
}

// UnimplementedChannelChunkServiceServer should be embedded to have forward compatible implementations.
type UnimplementedChannelChunkServiceServer struct {
}

func (UnimplementedChannelChunkServiceServer) CreateReplicas(ChannelChunkService_CreateReplicasServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateReplicas not implemented")
}
func (UnimplementedChannelChunkServiceServer) RetrieveReplicas(*ChannelChunkServiceRetrieveReplicasRequest, ChannelChunkService_RetrieveReplicasServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveReplicas not implemented")
}
func (UnimplementedChannelChunkServiceServer) DeleteReplicas(context.Context, *ChannelChunkServiceDeleteReplicasRequest) (*ChannelChunkServiceDeleteReplicasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteReplica not implemented")
}

// UnsafeChannelChunkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelChunkServiceServer will
// result in compilation errors.
type UnsafeChannelChunkServiceServer interface {
	mustEmbedUnimplementedChannelChunkServiceServer()
}

func RegisterChannelChunkServiceServer(s grpc.ServiceRegistrar, srv ChannelChunkServiceServer) {
	s.RegisterService(&ChannelChunkService_ServiceDesc, srv)
}

func _ChannelChunkService_CreateReplicas_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChannelChunkServiceServer).CreateReplicas(&channelChunkServiceCreateReplicasServer{stream})
}

type ChannelChunkService_CreateReplicasServer interface {
	SendAndClose(*ChannelChunkServiceCreateReplicasResponse) error
	Recv() (*ChannelChunkServiceCreateReplicasRequest, error)
	grpc.ServerStream
}

type channelChunkServiceCreateReplicasServer struct {
	grpc.ServerStream
}

func (x *channelChunkServiceCreateReplicasServer) SendAndClose(m *ChannelChunkServiceCreateReplicasResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *channelChunkServiceCreateReplicasServer) Recv() (*ChannelChunkServiceCreateReplicasRequest, error) {
	m := new(ChannelChunkServiceCreateReplicasRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChannelChunkService_RetrieveReplicas_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChannelChunkServiceRetrieveReplicasRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChannelChunkServiceServer).RetrieveReplicas(m, &channelChunkServiceRetrieveReplicasServer{stream})
}

type ChannelChunkService_RetrieveReplicasServer interface {
	Send(*ChannelChunkServiceRetrieveReplicasResponse) error
	grpc.ServerStream
}

type channelChunkServiceRetrieveReplicasServer struct {
	grpc.ServerStream
}

func (x *channelChunkServiceRetrieveReplicasServer) Send(m *ChannelChunkServiceRetrieveReplicasResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ChannelChunkService_DeleteReplicas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelChunkServiceDeleteReplicasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelChunkServiceServer).DeleteReplicas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chanchunk.v1.ChannelChunkService/DeleteReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelChunkServiceServer).DeleteReplicas(ctx, req.(*ChannelChunkServiceDeleteReplicasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelChunkService_ServiceDesc is the grpc.ServiceDesc for ChannelChunkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelChunkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chanchunk.v1.ChannelChunkService",
	HandlerType: (*ChannelChunkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteReplica",
			Handler:    _ChannelChunkService_DeleteReplicas_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateReplicas",
			Handler:       _ChannelChunkService_CreateReplicas_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RetrieveReplicas",
			Handler:       _ChannelChunkService_RetrieveReplicas_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chanchunk/v1/chanchunk.proto",
}
