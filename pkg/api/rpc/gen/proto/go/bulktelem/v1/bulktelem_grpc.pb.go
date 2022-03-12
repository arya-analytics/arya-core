// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package bulktelemv1

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

// BulkTelemServiceClient is the client API for BulkTelemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BulkTelemServiceClient interface {
	CreateStream(ctx context.Context, opts ...grpc.CallOption) (BulkTelemService_CreateStreamClient, error)
}

type bulkTelemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBulkTelemServiceClient(cc grpc.ClientConnInterface) BulkTelemServiceClient {
	return &bulkTelemServiceClient{cc}
}

func (c *bulkTelemServiceClient) CreateStream(ctx context.Context, opts ...grpc.CallOption) (BulkTelemService_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &BulkTelemService_ServiceDesc.Streams[0], "/bulktelem.v1.BulkTelemService/CreateStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &bulkTelemServiceCreateStreamClient{stream}
	return x, nil
}

type BulkTelemService_CreateStreamClient interface {
	Send(*CreateStreamRequest) error
	Recv() (*CreateStreamResponse, error)
	grpc.ClientStream
}

type bulkTelemServiceCreateStreamClient struct {
	grpc.ClientStream
}

func (x *bulkTelemServiceCreateStreamClient) Send(m *CreateStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *bulkTelemServiceCreateStreamClient) Recv() (*CreateStreamResponse, error) {
	m := new(CreateStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BulkTelemServiceServer is the server API for BulkTelemService service.
// All implementations should embed UnimplementedBulkTelemServiceServer
// for forward compatibility
type BulkTelemServiceServer interface {
	CreateStream(BulkTelemService_CreateStreamServer) error
}

// UnimplementedBulkTelemServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBulkTelemServiceServer struct {
}

func (UnimplementedBulkTelemServiceServer) CreateStream(BulkTelemService_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}

// UnsafeBulkTelemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BulkTelemServiceServer will
// result in compilation errors.
type UnsafeBulkTelemServiceServer interface {
	mustEmbedUnimplementedBulkTelemServiceServer()
}

func RegisterBulkTelemServiceServer(s grpc.ServiceRegistrar, srv BulkTelemServiceServer) {
	s.RegisterService(&BulkTelemService_ServiceDesc, srv)
}

func _BulkTelemService_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BulkTelemServiceServer).CreateStream(&bulkTelemServiceCreateStreamServer{stream})
}

type BulkTelemService_CreateStreamServer interface {
	Send(*CreateStreamResponse) error
	Recv() (*CreateStreamRequest, error)
	grpc.ServerStream
}

type bulkTelemServiceCreateStreamServer struct {
	grpc.ServerStream
}

func (x *bulkTelemServiceCreateStreamServer) Send(m *CreateStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *bulkTelemServiceCreateStreamServer) Recv() (*CreateStreamRequest, error) {
	m := new(CreateStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BulkTelemService_ServiceDesc is the grpc.ServiceDesc for BulkTelemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BulkTelemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bulktelem.v1.BulkTelemService",
	HandlerType: (*BulkTelemServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _BulkTelemService_CreateStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "bulktelem/v1/bulktelem.proto",
}
