// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/twitchtv/twirp-example/rpc/haberdasher/haberdasher.proto

/*
Package haberdasher is a generated protocol buffer package.

It is generated from these files:
	github.com/twitchtv/twirp-example/rpc/haberdasher/haberdasher.proto

It has these top-level messages:
	Hat
	Size
*/
package haberdasher

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

// A Hat is a piece of headwear made by a Haberdasher.
type Hat struct {
	// The size of a hat should always be in inches.
	Size int32 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
	// The color of a hat will never be 'invisible', but other than
	// that, anything is fair game.
	Color string `protobuf:"bytes,2,opt,name=color" json:"color,omitempty"`
	// The name of a hat is it's type. Like, 'bowler', or something.
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *Hat) Reset()                    { *m = Hat{} }
func (m *Hat) String() string            { return proto.CompactTextString(m) }
func (*Hat) ProtoMessage()               {}
func (*Hat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Hat) GetSize() int32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Hat) GetColor() string {
	if m != nil {
		return m.Color
	}
	return ""
}

func (m *Hat) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Size is passed when requesting a new hat to be made. It's always measured in
// inches.
type Size struct {
	Inches int32 `protobuf:"varint,1,opt,name=inches" json:"inches,omitempty"`
}

func (m *Size) Reset()                    { *m = Size{} }
func (m *Size) String() string            { return proto.CompactTextString(m) }
func (*Size) ProtoMessage()               {}
func (*Size) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Size) GetInches() int32 {
	if m != nil {
		return m.Inches
	}
	return 0
}

func init() {
	proto.RegisterType((*Hat)(nil), "twitch.twirp.example.haberdasher.Hat")
	proto.RegisterType((*Size)(nil), "twitch.twirp.example.haberdasher.Size")
}

func init() {
	proto.RegisterFile("github.com/twitchtv/twirp-example/rpc/haberdasher/haberdasher.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 214 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4e, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0x29, 0xcf, 0x2c, 0x49, 0xce, 0x28, 0x29, 0x03,
	0x31, 0x8a, 0x0a, 0x74, 0x53, 0x2b, 0x12, 0x73, 0x0b, 0x72, 0x52, 0xf5, 0x8b, 0x0a, 0x92, 0xf5,
	0x33, 0x12, 0x93, 0x52, 0x8b, 0x52, 0x12, 0x8b, 0x33, 0x52, 0x8b, 0x90, 0xd9, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x0a, 0x10, 0x9d, 0x7a, 0x60, 0x7d, 0x7a, 0x50, 0x7d, 0x7a, 0x48, 0xea,
	0x94, 0x9c, 0xb9, 0x98, 0x3d, 0x12, 0x4b, 0x84, 0x84, 0xb8, 0x58, 0x8a, 0x33, 0xab, 0x52, 0x25,
	0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0xc0, 0x6c, 0x21, 0x11, 0x2e, 0xd6, 0xe4, 0xfc, 0x9c, 0xfc,
	0x22, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x08, 0x07, 0xa4, 0x32, 0x2f, 0x31, 0x37, 0x55,
	0x82, 0x19, 0x2c, 0x08, 0x66, 0x2b, 0xc9, 0x71, 0xb1, 0x04, 0x83, 0x74, 0x88, 0x71, 0xb1, 0x65,
	0xe6, 0x25, 0x67, 0xa4, 0x16, 0x43, 0xcd, 0x81, 0xf2, 0x8c, 0xd2, 0xb9, 0xb8, 0x3d, 0x10, 0x76,
	0x0a, 0x45, 0x70, 0xb1, 0xfb, 0x26, 0x66, 0xa7, 0x82, 0xec, 0x55, 0xd3, 0x23, 0xe4, 0x42, 0x3d,
	0x90, 0xc9, 0x52, 0xaa, 0x84, 0xd5, 0x79, 0x24, 0x96, 0x38, 0xd9, 0x47, 0xd9, 0x92, 0x1c, 0x6c,
	0xd6, 0x48, 0xec, 0x24, 0x36, 0x70, 0xb8, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x73, 0xa6,
	0x91, 0xcf, 0x7e, 0x01, 0x00, 0x00,
}