// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: mq.proto

package mq_pb

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

// SeaweedMessagingClient is the client API for SeaweedMessaging service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeaweedMessagingClient interface {
	// control plane
	FindBrokerLeader(ctx context.Context, in *FindBrokerLeaderRequest, opts ...grpc.CallOption) (*FindBrokerLeaderResponse, error)
	AssignSegmentBrokers(ctx context.Context, in *AssignSegmentBrokersRequest, opts ...grpc.CallOption) (*AssignSegmentBrokersResponse, error)
	CheckSegmentStatus(ctx context.Context, in *CheckSegmentStatusRequest, opts ...grpc.CallOption) (*CheckSegmentStatusResponse, error)
	CheckBrokerLoad(ctx context.Context, in *CheckBrokerLoadRequest, opts ...grpc.CallOption) (*CheckBrokerLoadResponse, error)
	// control plane for balancer
	ConnectToBalancer(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_ConnectToBalancerClient, error)
	// control plane for topic partitions
	LookupTopicBrokers(ctx context.Context, in *LookupTopicBrokersRequest, opts ...grpc.CallOption) (*LookupTopicBrokersResponse, error)
	// a pub client will call this to get the topic partitions assignment
	RequestTopicPartitions(ctx context.Context, in *RequestTopicPartitionsRequest, opts ...grpc.CallOption) (*RequestTopicPartitionsResponse, error)
	AssignTopicPartitions(ctx context.Context, in *AssignTopicPartitionsRequest, opts ...grpc.CallOption) (*AssignTopicPartitionsResponse, error)
	CheckTopicPartitionsStatus(ctx context.Context, in *CheckTopicPartitionsStatusRequest, opts ...grpc.CallOption) (*CheckTopicPartitionsStatusResponse, error)
	// data plane
	Publish(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_PublishClient, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (SeaweedMessaging_SubscribeClient, error)
}

type seaweedMessagingClient struct {
	cc grpc.ClientConnInterface
}

func NewSeaweedMessagingClient(cc grpc.ClientConnInterface) SeaweedMessagingClient {
	return &seaweedMessagingClient{cc}
}

func (c *seaweedMessagingClient) FindBrokerLeader(ctx context.Context, in *FindBrokerLeaderRequest, opts ...grpc.CallOption) (*FindBrokerLeaderResponse, error) {
	out := new(FindBrokerLeaderResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/FindBrokerLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) AssignSegmentBrokers(ctx context.Context, in *AssignSegmentBrokersRequest, opts ...grpc.CallOption) (*AssignSegmentBrokersResponse, error) {
	out := new(AssignSegmentBrokersResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/AssignSegmentBrokers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) CheckSegmentStatus(ctx context.Context, in *CheckSegmentStatusRequest, opts ...grpc.CallOption) (*CheckSegmentStatusResponse, error) {
	out := new(CheckSegmentStatusResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/CheckSegmentStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) CheckBrokerLoad(ctx context.Context, in *CheckBrokerLoadRequest, opts ...grpc.CallOption) (*CheckBrokerLoadResponse, error) {
	out := new(CheckBrokerLoadResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/CheckBrokerLoad", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) ConnectToBalancer(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_ConnectToBalancerClient, error) {
	stream, err := c.cc.NewStream(ctx, &SeaweedMessaging_ServiceDesc.Streams[0], "/messaging_pb.SeaweedMessaging/ConnectToBalancer", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedMessagingConnectToBalancerClient{stream}
	return x, nil
}

type SeaweedMessaging_ConnectToBalancerClient interface {
	Send(*ConnectToBalancerRequest) error
	Recv() (*ConnectToBalancerResponse, error)
	grpc.ClientStream
}

type seaweedMessagingConnectToBalancerClient struct {
	grpc.ClientStream
}

func (x *seaweedMessagingConnectToBalancerClient) Send(m *ConnectToBalancerRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedMessagingConnectToBalancerClient) Recv() (*ConnectToBalancerResponse, error) {
	m := new(ConnectToBalancerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedMessagingClient) LookupTopicBrokers(ctx context.Context, in *LookupTopicBrokersRequest, opts ...grpc.CallOption) (*LookupTopicBrokersResponse, error) {
	out := new(LookupTopicBrokersResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/LookupTopicBrokers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) RequestTopicPartitions(ctx context.Context, in *RequestTopicPartitionsRequest, opts ...grpc.CallOption) (*RequestTopicPartitionsResponse, error) {
	out := new(RequestTopicPartitionsResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/RequestTopicPartitions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) AssignTopicPartitions(ctx context.Context, in *AssignTopicPartitionsRequest, opts ...grpc.CallOption) (*AssignTopicPartitionsResponse, error) {
	out := new(AssignTopicPartitionsResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/AssignTopicPartitions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) CheckTopicPartitionsStatus(ctx context.Context, in *CheckTopicPartitionsStatusRequest, opts ...grpc.CallOption) (*CheckTopicPartitionsStatusResponse, error) {
	out := new(CheckTopicPartitionsStatusResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/CheckTopicPartitionsStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) Publish(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_PublishClient, error) {
	stream, err := c.cc.NewStream(ctx, &SeaweedMessaging_ServiceDesc.Streams[1], "/messaging_pb.SeaweedMessaging/Publish", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedMessagingPublishClient{stream}
	return x, nil
}

type SeaweedMessaging_PublishClient interface {
	Send(*PublishRequest) error
	Recv() (*PublishResponse, error)
	grpc.ClientStream
}

type seaweedMessagingPublishClient struct {
	grpc.ClientStream
}

func (x *seaweedMessagingPublishClient) Send(m *PublishRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedMessagingPublishClient) Recv() (*PublishResponse, error) {
	m := new(PublishResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedMessagingClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (SeaweedMessaging_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &SeaweedMessaging_ServiceDesc.Streams[2], "/messaging_pb.SeaweedMessaging/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedMessagingSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SeaweedMessaging_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type seaweedMessagingSubscribeClient struct {
	grpc.ClientStream
}

func (x *seaweedMessagingSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SeaweedMessagingServer is the server API for SeaweedMessaging service.
// All implementations must embed UnimplementedSeaweedMessagingServer
// for forward compatibility
type SeaweedMessagingServer interface {
	// control plane
	FindBrokerLeader(context.Context, *FindBrokerLeaderRequest) (*FindBrokerLeaderResponse, error)
	AssignSegmentBrokers(context.Context, *AssignSegmentBrokersRequest) (*AssignSegmentBrokersResponse, error)
	CheckSegmentStatus(context.Context, *CheckSegmentStatusRequest) (*CheckSegmentStatusResponse, error)
	CheckBrokerLoad(context.Context, *CheckBrokerLoadRequest) (*CheckBrokerLoadResponse, error)
	// control plane for balancer
	ConnectToBalancer(SeaweedMessaging_ConnectToBalancerServer) error
	// control plane for topic partitions
	LookupTopicBrokers(context.Context, *LookupTopicBrokersRequest) (*LookupTopicBrokersResponse, error)
	// a pub client will call this to get the topic partitions assignment
	RequestTopicPartitions(context.Context, *RequestTopicPartitionsRequest) (*RequestTopicPartitionsResponse, error)
	AssignTopicPartitions(context.Context, *AssignTopicPartitionsRequest) (*AssignTopicPartitionsResponse, error)
	CheckTopicPartitionsStatus(context.Context, *CheckTopicPartitionsStatusRequest) (*CheckTopicPartitionsStatusResponse, error)
	// data plane
	Publish(SeaweedMessaging_PublishServer) error
	Subscribe(*SubscribeRequest, SeaweedMessaging_SubscribeServer) error
	mustEmbedUnimplementedSeaweedMessagingServer()
}

// UnimplementedSeaweedMessagingServer must be embedded to have forward compatible implementations.
type UnimplementedSeaweedMessagingServer struct {
}

func (UnimplementedSeaweedMessagingServer) FindBrokerLeader(context.Context, *FindBrokerLeaderRequest) (*FindBrokerLeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindBrokerLeader not implemented")
}
func (UnimplementedSeaweedMessagingServer) AssignSegmentBrokers(context.Context, *AssignSegmentBrokersRequest) (*AssignSegmentBrokersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignSegmentBrokers not implemented")
}
func (UnimplementedSeaweedMessagingServer) CheckSegmentStatus(context.Context, *CheckSegmentStatusRequest) (*CheckSegmentStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSegmentStatus not implemented")
}
func (UnimplementedSeaweedMessagingServer) CheckBrokerLoad(context.Context, *CheckBrokerLoadRequest) (*CheckBrokerLoadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckBrokerLoad not implemented")
}
func (UnimplementedSeaweedMessagingServer) ConnectToBalancer(SeaweedMessaging_ConnectToBalancerServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectToBalancer not implemented")
}
func (UnimplementedSeaweedMessagingServer) LookupTopicBrokers(context.Context, *LookupTopicBrokersRequest) (*LookupTopicBrokersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LookupTopicBrokers not implemented")
}
func (UnimplementedSeaweedMessagingServer) RequestTopicPartitions(context.Context, *RequestTopicPartitionsRequest) (*RequestTopicPartitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestTopicPartitions not implemented")
}
func (UnimplementedSeaweedMessagingServer) AssignTopicPartitions(context.Context, *AssignTopicPartitionsRequest) (*AssignTopicPartitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignTopicPartitions not implemented")
}
func (UnimplementedSeaweedMessagingServer) CheckTopicPartitionsStatus(context.Context, *CheckTopicPartitionsStatusRequest) (*CheckTopicPartitionsStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckTopicPartitionsStatus not implemented")
}
func (UnimplementedSeaweedMessagingServer) Publish(SeaweedMessaging_PublishServer) error {
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedSeaweedMessagingServer) Subscribe(*SubscribeRequest, SeaweedMessaging_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedSeaweedMessagingServer) mustEmbedUnimplementedSeaweedMessagingServer() {}

// UnsafeSeaweedMessagingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeaweedMessagingServer will
// result in compilation errors.
type UnsafeSeaweedMessagingServer interface {
	mustEmbedUnimplementedSeaweedMessagingServer()
}

func RegisterSeaweedMessagingServer(s grpc.ServiceRegistrar, srv SeaweedMessagingServer) {
	s.RegisterService(&SeaweedMessaging_ServiceDesc, srv)
}

func _SeaweedMessaging_FindBrokerLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindBrokerLeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).FindBrokerLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/FindBrokerLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).FindBrokerLeader(ctx, req.(*FindBrokerLeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_AssignSegmentBrokers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignSegmentBrokersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).AssignSegmentBrokers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/AssignSegmentBrokers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).AssignSegmentBrokers(ctx, req.(*AssignSegmentBrokersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_CheckSegmentStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckSegmentStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).CheckSegmentStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/CheckSegmentStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).CheckSegmentStatus(ctx, req.(*CheckSegmentStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_CheckBrokerLoad_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckBrokerLoadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).CheckBrokerLoad(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/CheckBrokerLoad",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).CheckBrokerLoad(ctx, req.(*CheckBrokerLoadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_ConnectToBalancer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedMessagingServer).ConnectToBalancer(&seaweedMessagingConnectToBalancerServer{stream})
}

type SeaweedMessaging_ConnectToBalancerServer interface {
	Send(*ConnectToBalancerResponse) error
	Recv() (*ConnectToBalancerRequest, error)
	grpc.ServerStream
}

type seaweedMessagingConnectToBalancerServer struct {
	grpc.ServerStream
}

func (x *seaweedMessagingConnectToBalancerServer) Send(m *ConnectToBalancerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedMessagingConnectToBalancerServer) Recv() (*ConnectToBalancerRequest, error) {
	m := new(ConnectToBalancerRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SeaweedMessaging_LookupTopicBrokers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupTopicBrokersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).LookupTopicBrokers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/LookupTopicBrokers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).LookupTopicBrokers(ctx, req.(*LookupTopicBrokersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_RequestTopicPartitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestTopicPartitionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).RequestTopicPartitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/RequestTopicPartitions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).RequestTopicPartitions(ctx, req.(*RequestTopicPartitionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_AssignTopicPartitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignTopicPartitionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).AssignTopicPartitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/AssignTopicPartitions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).AssignTopicPartitions(ctx, req.(*AssignTopicPartitionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_CheckTopicPartitionsStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckTopicPartitionsStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).CheckTopicPartitionsStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/CheckTopicPartitionsStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).CheckTopicPartitionsStatus(ctx, req.(*CheckTopicPartitionsStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_Publish_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedMessagingServer).Publish(&seaweedMessagingPublishServer{stream})
}

type SeaweedMessaging_PublishServer interface {
	Send(*PublishResponse) error
	Recv() (*PublishRequest, error)
	grpc.ServerStream
}

type seaweedMessagingPublishServer struct {
	grpc.ServerStream
}

func (x *seaweedMessagingPublishServer) Send(m *PublishResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedMessagingPublishServer) Recv() (*PublishRequest, error) {
	m := new(PublishRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SeaweedMessaging_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SeaweedMessagingServer).Subscribe(m, &seaweedMessagingSubscribeServer{stream})
}

type SeaweedMessaging_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type seaweedMessagingSubscribeServer struct {
	grpc.ServerStream
}

func (x *seaweedMessagingSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// SeaweedMessaging_ServiceDesc is the grpc.ServiceDesc for SeaweedMessaging service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SeaweedMessaging_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messaging_pb.SeaweedMessaging",
	HandlerType: (*SeaweedMessagingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindBrokerLeader",
			Handler:    _SeaweedMessaging_FindBrokerLeader_Handler,
		},
		{
			MethodName: "AssignSegmentBrokers",
			Handler:    _SeaweedMessaging_AssignSegmentBrokers_Handler,
		},
		{
			MethodName: "CheckSegmentStatus",
			Handler:    _SeaweedMessaging_CheckSegmentStatus_Handler,
		},
		{
			MethodName: "CheckBrokerLoad",
			Handler:    _SeaweedMessaging_CheckBrokerLoad_Handler,
		},
		{
			MethodName: "LookupTopicBrokers",
			Handler:    _SeaweedMessaging_LookupTopicBrokers_Handler,
		},
		{
			MethodName: "RequestTopicPartitions",
			Handler:    _SeaweedMessaging_RequestTopicPartitions_Handler,
		},
		{
			MethodName: "AssignTopicPartitions",
			Handler:    _SeaweedMessaging_AssignTopicPartitions_Handler,
		},
		{
			MethodName: "CheckTopicPartitionsStatus",
			Handler:    _SeaweedMessaging_CheckTopicPartitionsStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectToBalancer",
			Handler:       _SeaweedMessaging_ConnectToBalancer_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Publish",
			Handler:       _SeaweedMessaging_Publish_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _SeaweedMessaging_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mq.proto",
}
