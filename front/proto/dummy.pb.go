// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dummy.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	dummy.proto

It has these top-level messages:
	DummyRequest
	DummyResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type DummyRequest struct {
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *DummyRequest) Reset()                    { *m = DummyRequest{} }
func (m *DummyRequest) String() string            { return proto1.CompactTextString(m) }
func (*DummyRequest) ProtoMessage()               {}
func (*DummyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DummyRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type DummyResponse struct {
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *DummyResponse) Reset()                    { *m = DummyResponse{} }
func (m *DummyResponse) String() string            { return proto1.CompactTextString(m) }
func (*DummyResponse) ProtoMessage()               {}
func (*DummyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DummyResponse) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto1.RegisterType((*DummyRequest)(nil), "proto.DummyRequest")
	proto1.RegisterType((*DummyResponse)(nil), "proto.DummyResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GetDummyData service

type GetDummyDataClient interface {
	Get(ctx context.Context, in *DummyRequest, opts ...grpc.CallOption) (*DummyResponse, error)
}

type getDummyDataClient struct {
	cc *grpc.ClientConn
}

func NewGetDummyDataClient(cc *grpc.ClientConn) GetDummyDataClient {
	return &getDummyDataClient{cc}
}

func (c *getDummyDataClient) Get(ctx context.Context, in *DummyRequest, opts ...grpc.CallOption) (*DummyResponse, error) {
	out := new(DummyResponse)
	err := grpc.Invoke(ctx, "/proto.GetDummyData/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GetDummyData service

type GetDummyDataServer interface {
	Get(context.Context, *DummyRequest) (*DummyResponse, error)
}

func RegisterGetDummyDataServer(s *grpc.Server, srv GetDummyDataServer) {
	s.RegisterService(&_GetDummyData_serviceDesc, srv)
}

func _GetDummyData_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DummyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetDummyDataServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GetDummyData/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetDummyDataServer).Get(ctx, req.(*DummyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GetDummyData_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GetDummyData",
	HandlerType: (*GetDummyDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _GetDummyData_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dummy.proto",
}

func init() { proto1.RegisterFile("dummy.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x29, 0xcd, 0xcd,
	0xad, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x1a, 0x5c, 0x3c, 0x2e,
	0x20, 0xd1, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x09, 0x2e, 0xf6, 0x82, 0xc4, 0xca,
	0x9c, 0xfc, 0xc4, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x18, 0x57, 0x49, 0x93, 0x8b,
	0x17, 0xaa, 0xb2, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x15, 0xb7, 0x52, 0x23, 0x07, 0x2e, 0x1e, 0xf7,
	0xd4, 0x12, 0xb0, 0x6a, 0x97, 0xc4, 0x92, 0x44, 0x21, 0x03, 0x2e, 0x66, 0xf7, 0xd4, 0x12, 0x21,
	0x61, 0x88, 0xd5, 0x7a, 0xc8, 0x16, 0x4a, 0x89, 0xa0, 0x0a, 0x42, 0xcc, 0x4e, 0x62, 0x03, 0x0b,
	0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb4, 0xe0, 0x9a, 0x91, 0xb3, 0x00, 0x00, 0x00,
}
