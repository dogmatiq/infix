// Code generated by protoc-gen-go. DO NOT EDIT.
// source: draftspecs/messagingspec/error.proto

package messagingspec

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// UnrecognizedApplication is an error-details value for INVALID_ARGUMENT
// errors that occurred because a specific application key was not recognized by
// the server.
type UnrecognizedApplication struct {
	// ApplicationKey is the identity of the application that produced the error.
	ApplicationKey       string   `protobuf:"bytes,1,opt,name=application_key,json=applicationKey,proto3" json:"application_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnrecognizedApplication) Reset()         { *m = UnrecognizedApplication{} }
func (m *UnrecognizedApplication) String() string { return proto.CompactTextString(m) }
func (*UnrecognizedApplication) ProtoMessage()    {}
func (*UnrecognizedApplication) Descriptor() ([]byte, []int) {
	return fileDescriptor_8112ac07a2bc8da6, []int{0}
}

func (m *UnrecognizedApplication) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnrecognizedApplication.Unmarshal(m, b)
}
func (m *UnrecognizedApplication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnrecognizedApplication.Marshal(b, m, deterministic)
}
func (m *UnrecognizedApplication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnrecognizedApplication.Merge(m, src)
}
func (m *UnrecognizedApplication) XXX_Size() int {
	return xxx_messageInfo_UnrecognizedApplication.Size(m)
}
func (m *UnrecognizedApplication) XXX_DiscardUnknown() {
	xxx_messageInfo_UnrecognizedApplication.DiscardUnknown(m)
}

var xxx_messageInfo_UnrecognizedApplication proto.InternalMessageInfo

func (m *UnrecognizedApplication) GetApplicationKey() string {
	if m != nil {
		return m.ApplicationKey
	}
	return ""
}

// UnrecognizedMessage is an error-details value for INVALID_ARGUMENT errors
// that occurred because a specific message type was not recognized by the
// server.
type UnrecognizedMessage struct {
	// Name is the name of the message type.
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnrecognizedMessage) Reset()         { *m = UnrecognizedMessage{} }
func (m *UnrecognizedMessage) String() string { return proto.CompactTextString(m) }
func (*UnrecognizedMessage) ProtoMessage()    {}
func (*UnrecognizedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8112ac07a2bc8da6, []int{1}
}

func (m *UnrecognizedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnrecognizedMessage.Unmarshal(m, b)
}
func (m *UnrecognizedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnrecognizedMessage.Marshal(b, m, deterministic)
}
func (m *UnrecognizedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnrecognizedMessage.Merge(m, src)
}
func (m *UnrecognizedMessage) XXX_Size() int {
	return xxx_messageInfo_UnrecognizedMessage.Size(m)
}
func (m *UnrecognizedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_UnrecognizedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_UnrecognizedMessage proto.InternalMessageInfo

func (m *UnrecognizedMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*UnrecognizedApplication)(nil), "dogma.messaging.v1.UnrecognizedApplication")
	proto.RegisterType((*UnrecognizedMessage)(nil), "dogma.messaging.v1.UnrecognizedMessage")
}

func init() {
	proto.RegisterFile("draftspecs/messagingspec/error.proto", fileDescriptor_8112ac07a2bc8da6)
}

var fileDescriptor_8112ac07a2bc8da6 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0x29, 0x4a, 0x4c,
	0x2b, 0x29, 0x2e, 0x48, 0x4d, 0x2e, 0xd6, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0xcf, 0xcc, 0x4b,
	0x07, 0x71, 0xf5, 0x53, 0x8b, 0x8a, 0xf2, 0x8b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x84,
	0x52, 0xf2, 0xd3, 0x73, 0x13, 0xf5, 0xe0, 0x0a, 0xf4, 0xca, 0x0c, 0x95, 0x9c, 0xb8, 0xc4, 0x43,
	0xf3, 0x8a, 0x52, 0x93, 0xf3, 0xd3, 0xf3, 0x32, 0xab, 0x52, 0x53, 0x1c, 0x0b, 0x0a, 0x72, 0x32,
	0x93, 0x13, 0x4b, 0x32, 0xf3, 0xf3, 0x84, 0xd4, 0xb9, 0xf8, 0x13, 0x11, 0xdc, 0xf8, 0xec, 0xd4,
	0x4a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x3e, 0x24, 0x61, 0xef, 0xd4, 0x4a, 0x25, 0x4d,
	0x2e, 0x61, 0x64, 0x33, 0x7c, 0xc1, 0xe6, 0xa7, 0x0a, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6,
	0x4a, 0x30, 0x81, 0x35, 0x81, 0xd9, 0x4e, 0x26, 0x51, 0x46, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49,
	0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x60, 0xf7, 0x94, 0x64, 0x16, 0xea, 0x67, 0xe6, 0xa5, 0x65, 0x56,
	0xe8, 0xe3, 0xf2, 0x44, 0x12, 0x1b, 0xd8, 0xfd, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x55,
	0x1e, 0x45, 0x7a, 0xe7, 0x00, 0x00, 0x00,
}
