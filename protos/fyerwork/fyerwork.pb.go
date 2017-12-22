// Code generated by protoc-gen-go.
// source: protos/fyerwork/fyerwork.proto
// DO NOT EDIT!

/*
Package fyerwork is a generated protocol buffer package.

It is generated from these files:
	protos/fyerwork/fyerwork.proto

It has these top-level messages:
	FetchRequest
	FetchResponse
*/
package fyerwork

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FetchRequest struct {
	Name  string              `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Range *FetchRequest_Range `protobuf:"bytes,2,opt,name=range" json:"range,omitempty"`
}

func (m *FetchRequest) Reset()                    { *m = FetchRequest{} }
func (m *FetchRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchRequest) ProtoMessage()               {}
func (*FetchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FetchRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FetchRequest) GetRange() *FetchRequest_Range {
	if m != nil {
		return m.Range
	}
	return nil
}

type FetchRequest_Range struct {
	Start int64 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	Size  int64 `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
}

func (m *FetchRequest_Range) Reset()                    { *m = FetchRequest_Range{} }
func (m *FetchRequest_Range) String() string            { return proto.CompactTextString(m) }
func (*FetchRequest_Range) ProtoMessage()               {}
func (*FetchRequest_Range) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *FetchRequest_Range) GetStart() int64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *FetchRequest_Range) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type FetchResponse struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *FetchResponse) Reset()                    { *m = FetchResponse{} }
func (m *FetchResponse) String() string            { return proto.CompactTextString(m) }
func (*FetchResponse) ProtoMessage()               {}
func (*FetchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *FetchResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*FetchRequest)(nil), "cn.zumium.fyerwork.FetchRequest")
	proto.RegisterType((*FetchRequest_Range)(nil), "cn.zumium.fyerwork.FetchRequest.Range")
	proto.RegisterType((*FetchResponse)(nil), "cn.zumium.fyerwork.FetchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Fyerwork service

type FyerworkClient interface {
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error)
}

type fyerworkClient struct {
	cc *grpc.ClientConn
}

func NewFyerworkClient(cc *grpc.ClientConn) FyerworkClient {
	return &fyerworkClient{cc}
}

func (c *fyerworkClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error) {
	out := new(FetchResponse)
	err := grpc.Invoke(ctx, "/cn.zumium.fyerwork.Fyerwork/Fetch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Fyerwork service

type FyerworkServer interface {
	Fetch(context.Context, *FetchRequest) (*FetchResponse, error)
}

func RegisterFyerworkServer(s *grpc.Server, srv FyerworkServer) {
	s.RegisterService(&_Fyerwork_serviceDesc, srv)
}

func _Fyerwork_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FyerworkServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cn.zumium.fyerwork.Fyerwork/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FyerworkServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Fyerwork_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cn.zumium.fyerwork.Fyerwork",
	HandlerType: (*FyerworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _Fyerwork_Fetch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/fyerwork/fyerwork.proto",
}

func init() { proto.RegisterFile("protos/fyerwork/fyerwork.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x4f, 0xab, 0x4c, 0x2d, 0x2a, 0xcf, 0x2f, 0xca, 0x86, 0x33, 0xf4, 0xc0, 0x12,
	0x42, 0x42, 0xc9, 0x79, 0x7a, 0x55, 0xa5, 0xb9, 0x99, 0xa5, 0xb9, 0x7a, 0x30, 0x19, 0xa5, 0xc9,
	0x8c, 0x5c, 0x3c, 0x6e, 0xa9, 0x25, 0xc9, 0x19, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42,
	0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6,
	0x90, 0x0d, 0x17, 0x6b, 0x51, 0x62, 0x5e, 0x7a, 0xaa, 0x04, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91,
	0x9a, 0x1e, 0xa6, 0x41, 0x7a, 0xc8, 0x86, 0xe8, 0x05, 0x81, 0x54, 0x07, 0x41, 0x34, 0x49, 0x19,
	0x72, 0xb1, 0x82, 0xf9, 0x42, 0x22, 0x5c, 0xac, 0xc5, 0x25, 0x89, 0x45, 0x25, 0x60, 0xb3, 0x99,
	0x83, 0x20, 0x1c, 0x90, 0x85, 0xc5, 0x99, 0x55, 0x10, 0xb3, 0x99, 0x83, 0xc0, 0x6c, 0x25, 0x65,
	0x2e, 0x5e, 0xa8, 0x79, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x20, 0x45, 0x29, 0x89, 0x25, 0x89,
	0x60, 0x9d, 0x3c, 0x41, 0x60, 0xb6, 0x51, 0x04, 0x17, 0x87, 0x1b, 0xd4, 0x76, 0x21, 0x1f, 0x2e,
	0x56, 0xb0, 0x06, 0x21, 0x05, 0x42, 0x6e, 0x93, 0x52, 0xc4, 0xa3, 0x02, 0x62, 0x9b, 0x93, 0x3e,
	0x97, 0x22, 0x16, 0x35, 0xe0, 0x40, 0x84, 0x73, 0x03, 0x18, 0xa3, 0x38, 0x60, 0xec, 0x24, 0x36,
	0xb0, 0x9c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x35, 0x82, 0x68, 0x02, 0x82, 0x01, 0x00, 0x00,
}
