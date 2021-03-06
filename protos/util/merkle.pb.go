// Code generated by protoc-gen-go.
// source: protos/util/merkle.proto
// DO NOT EDIT!

/*
Package util is a generated protocol buffer package.

It is generated from these files:
	protos/util/merkle.proto

It has these top-level messages:
	MerkleItem
	MerklePack
*/
package util

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MerkleItem_Type int32

const (
	MerkleItem_EMPTY MerkleItem_Type = 0
	MerkleItem_DATA  MerkleItem_Type = 1
)

var MerkleItem_Type_name = map[int32]string{
	0: "EMPTY",
	1: "DATA",
}
var MerkleItem_Type_value = map[string]int32{
	"EMPTY": 0,
	"DATA":  1,
}

func (x MerkleItem_Type) String() string {
	return proto.EnumName(MerkleItem_Type_name, int32(x))
}
func (MerkleItem_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type MerkleItem struct {
	Type MerkleItem_Type `protobuf:"varint,1,opt,name=type,enum=cn.zumium.fyer.util.MerkleItem_Type" json:"type,omitempty"`
	Data []byte          `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *MerkleItem) Reset()                    { *m = MerkleItem{} }
func (m *MerkleItem) String() string            { return proto.CompactTextString(m) }
func (*MerkleItem) ProtoMessage()               {}
func (*MerkleItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MerkleItem) GetType() MerkleItem_Type {
	if m != nil {
		return m.Type
	}
	return MerkleItem_EMPTY
}

func (m *MerkleItem) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type MerklePack struct {
	Size  uint64        `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
	Items []*MerkleItem `protobuf:"bytes,2,rep,name=items" json:"items,omitempty"`
}

func (m *MerklePack) Reset()                    { *m = MerklePack{} }
func (m *MerklePack) String() string            { return proto.CompactTextString(m) }
func (*MerklePack) ProtoMessage()               {}
func (*MerklePack) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MerklePack) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *MerklePack) GetItems() []*MerkleItem {
	if m != nil {
		return m.Items
	}
	return nil
}

func init() {
	proto.RegisterType((*MerkleItem)(nil), "cn.zumium.fyer.util.MerkleItem")
	proto.RegisterType((*MerklePack)(nil), "cn.zumium.fyer.util.MerklePack")
	proto.RegisterEnum("cn.zumium.fyer.util.MerkleItem_Type", MerkleItem_Type_name, MerkleItem_Type_value)
}

func init() { proto.RegisterFile("protos/util/merkle.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x2f, 0x2d, 0xc9, 0xcc, 0xd1, 0xcf, 0x4d, 0x2d, 0xca, 0xce, 0x49, 0xd5, 0x03,
	0x0b, 0x09, 0x09, 0x27, 0xe7, 0xe9, 0x55, 0x95, 0xe6, 0x66, 0x96, 0xe6, 0xea, 0xa5, 0x55, 0xa6,
	0x16, 0xe9, 0x81, 0x54, 0x28, 0x95, 0x73, 0x71, 0xf9, 0x82, 0x15, 0x79, 0x96, 0xa4, 0xe6, 0x0a,
	0x59, 0x70, 0xb1, 0x94, 0x54, 0x16, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x19, 0xa9, 0xe8,
	0x61, 0xd1, 0xa1, 0x87, 0x50, 0xae, 0x17, 0x52, 0x59, 0x90, 0x1a, 0x04, 0xd6, 0x21, 0x24, 0xc4,
	0xc5, 0x92, 0x92, 0x58, 0x92, 0x28, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x13, 0x04, 0x66, 0x2b, 0x49,
	0x73, 0xb1, 0x80, 0x54, 0x08, 0x71, 0x72, 0xb1, 0xba, 0xfa, 0x06, 0x84, 0x44, 0x0a, 0x30, 0x08,
	0x71, 0x70, 0xb1, 0xb8, 0x38, 0x86, 0x38, 0x0a, 0x30, 0x2a, 0x85, 0xc3, 0x2c, 0x0e, 0x48, 0x4c,
	0xce, 0x06, 0x69, 0x2f, 0xce, 0xac, 0x82, 0x58, 0xcc, 0x12, 0x04, 0x66, 0x0b, 0x99, 0x72, 0xb1,
	0x66, 0x96, 0xa4, 0xe6, 0x16, 0x4b, 0x30, 0x29, 0x30, 0x6b, 0x70, 0x1b, 0xc9, 0x13, 0x70, 0x4d,
	0x10, 0x44, 0xb5, 0x93, 0x22, 0x97, 0x24, 0x9a, 0x42, 0xb0, 0xf7, 0xc1, 0xca, 0xa3, 0x58, 0x40,
	0x64, 0x12, 0x1b, 0x58, 0xc4, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x72, 0x0f, 0x65, 0xa6, 0x2c,
	0x01, 0x00, 0x00,
}
