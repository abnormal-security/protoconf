// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protoconfmutation

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

// ProtoconfMutationServiceClient is the client API for ProtoconfMutationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProtoconfMutationServiceClient interface {
	MutateConfig(ctx context.Context, in *ConfigMutationRequest, opts ...grpc.CallOption) (*ConfigMutationResponse, error)
}

type protoconfMutationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProtoconfMutationServiceClient(cc grpc.ClientConnInterface) ProtoconfMutationServiceClient {
	return &protoconfMutationServiceClient{cc}
}

func (c *protoconfMutationServiceClient) MutateConfig(ctx context.Context, in *ConfigMutationRequest, opts ...grpc.CallOption) (*ConfigMutationResponse, error) {
	out := new(ConfigMutationResponse)
	err := c.cc.Invoke(ctx, "/v1.ProtoconfMutationService/MutateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProtoconfMutationServiceServer is the server API for ProtoconfMutationService service.
// All implementations must embed UnimplementedProtoconfMutationServiceServer
// for forward compatibility
type ProtoconfMutationServiceServer interface {
	MutateConfig(context.Context, *ConfigMutationRequest) (*ConfigMutationResponse, error)
	mustEmbedUnimplementedProtoconfMutationServiceServer()
}

// UnimplementedProtoconfMutationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProtoconfMutationServiceServer struct {
}

func (UnimplementedProtoconfMutationServiceServer) MutateConfig(context.Context, *ConfigMutationRequest) (*ConfigMutationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MutateConfig not implemented")
}
func (UnimplementedProtoconfMutationServiceServer) mustEmbedUnimplementedProtoconfMutationServiceServer() {
}

// UnsafeProtoconfMutationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProtoconfMutationServiceServer will
// result in compilation errors.
type UnsafeProtoconfMutationServiceServer interface {
	mustEmbedUnimplementedProtoconfMutationServiceServer()
}

func RegisterProtoconfMutationServiceServer(s grpc.ServiceRegistrar, srv ProtoconfMutationServiceServer) {
	s.RegisterService(&ProtoconfMutationService_ServiceDesc, srv)
}

func _ProtoconfMutationService_MutateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigMutationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoconfMutationServiceServer).MutateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ProtoconfMutationService/MutateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoconfMutationServiceServer).MutateConfig(ctx, req.(*ConfigMutationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProtoconfMutationService_ServiceDesc is the grpc.ServiceDesc for ProtoconfMutationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProtoconfMutationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ProtoconfMutationService",
	HandlerType: (*ProtoconfMutationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MutateConfig",
			Handler:    _ProtoconfMutationService_MutateConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protoconf_mutation.proto",
}