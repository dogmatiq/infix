// Code generated by protoc-gen-go. DO NOT EDIT.
// source: draftspecs/messagingspec/messagetypes.proto

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

type MessageTypesRequest struct {
	// ApplicationKey is the identity key of the application to query.
	ApplicationKey       string   `protobuf:"bytes,1,opt,name=application_key,json=applicationKey,proto3" json:"application_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageTypesRequest) Reset()         { *m = MessageTypesRequest{} }
func (m *MessageTypesRequest) String() string { return proto.CompactTextString(m) }
func (*MessageTypesRequest) ProtoMessage()    {}
func (*MessageTypesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b2a3c59525cd061, []int{0}
}

func (m *MessageTypesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageTypesRequest.Unmarshal(m, b)
}
func (m *MessageTypesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageTypesRequest.Marshal(b, m, deterministic)
}
func (m *MessageTypesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageTypesRequest.Merge(m, src)
}
func (m *MessageTypesRequest) XXX_Size() int {
	return xxx_messageInfo_MessageTypesRequest.Size(m)
}
func (m *MessageTypesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageTypesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageTypesRequest proto.InternalMessageInfo

func (m *MessageTypesRequest) GetApplicationKey() string {
	if m != nil {
		return m.ApplicationKey
	}
	return ""
}

type MessageTypesResponse struct {
	// MessageTypes is the set of messages supported by the server.
	MessageTypes         []*MessageType `protobuf:"bytes,1,rep,name=message_types,json=messageTypes,proto3" json:"message_types,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *MessageTypesResponse) Reset()         { *m = MessageTypesResponse{} }
func (m *MessageTypesResponse) String() string { return proto.CompactTextString(m) }
func (*MessageTypesResponse) ProtoMessage()    {}
func (*MessageTypesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b2a3c59525cd061, []int{1}
}

func (m *MessageTypesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageTypesResponse.Unmarshal(m, b)
}
func (m *MessageTypesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageTypesResponse.Marshal(b, m, deterministic)
}
func (m *MessageTypesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageTypesResponse.Merge(m, src)
}
func (m *MessageTypesResponse) XXX_Size() int {
	return xxx_messageInfo_MessageTypesResponse.Size(m)
}
func (m *MessageTypesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageTypesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MessageTypesResponse proto.InternalMessageInfo

func (m *MessageTypesResponse) GetMessageTypes() []*MessageType {
	if m != nil {
		return m.MessageTypes
	}
	return nil
}

type MessageType struct {
	// PortableName is the unique name used to identify messages of this type.
	PortableName string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	// ConfigName is the name used to identify this message type in the
	// dogma.config.v1 API.
	//
	// This name may differ across builds, as it is based on the fully-qualified
	// Go type name.
	ConfigName string `protobuf:"bytes,2,opt,name=config_name,json=configName,proto3" json:"config_name,omitempty"`
	// MediaTypes is a list of MIME media-types that the server may use to
	// represent messages of this type.
	MediaTypes           []string `protobuf:"bytes,3,rep,name=media_types,json=mediaTypes,proto3" json:"media_types,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageType) Reset()         { *m = MessageType{} }
func (m *MessageType) String() string { return proto.CompactTextString(m) }
func (*MessageType) ProtoMessage()    {}
func (*MessageType) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b2a3c59525cd061, []int{2}
}

func (m *MessageType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageType.Unmarshal(m, b)
}
func (m *MessageType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageType.Marshal(b, m, deterministic)
}
func (m *MessageType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageType.Merge(m, src)
}
func (m *MessageType) XXX_Size() int {
	return xxx_messageInfo_MessageType.Size(m)
}
func (m *MessageType) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageType.DiscardUnknown(m)
}

var xxx_messageInfo_MessageType proto.InternalMessageInfo

func (m *MessageType) GetPortableName() string {
	if m != nil {
		return m.PortableName
	}
	return ""
}

func (m *MessageType) GetConfigName() string {
	if m != nil {
		return m.ConfigName
	}
	return ""
}

func (m *MessageType) GetMediaTypes() []string {
	if m != nil {
		return m.MediaTypes
	}
	return nil
}

func init() {
	proto.RegisterType((*MessageTypesRequest)(nil), "dogma.messaging.v1.MessageTypesRequest")
	proto.RegisterType((*MessageTypesResponse)(nil), "dogma.messaging.v1.MessageTypesResponse")
	proto.RegisterType((*MessageType)(nil), "dogma.messaging.v1.MessageType")
}

func init() {
	proto.RegisterFile("draftspecs/messagingspec/messagetypes.proto", fileDescriptor_0b2a3c59525cd061)
}

var fileDescriptor_0b2a3c59525cd061 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xcf, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x89, 0x01, 0xc1, 0x4d, 0xab, 0xb0, 0x7a, 0xe8, 0xad, 0x21, 0x1e, 0x2c, 0x08, 0x1b,
	0xac, 0x9e, 0x3d, 0x88, 0x37, 0xd1, 0x43, 0xf0, 0x24, 0x42, 0xd8, 0x24, 0x93, 0x38, 0xd8, 0xfd,
	0xd1, 0xec, 0x54, 0xcc, 0x7f, 0x2f, 0xd9, 0xa4, 0x12, 0x91, 0x1e, 0xe7, 0xed, 0xb7, 0x33, 0x1f,
	0x8f, 0x5d, 0x57, 0xad, 0xac, 0xc9, 0x59, 0x28, 0x5d, 0xaa, 0xc0, 0x39, 0xd9, 0xa0, 0x6e, 0xfa,
	0x71, 0x9c, 0x80, 0x3a, 0x0b, 0x4e, 0xd8, 0xd6, 0x90, 0xe1, 0xbc, 0x32, 0x8d, 0x92, 0xe2, 0x97,
	0x13, 0x5f, 0x37, 0xc9, 0x3d, 0x3b, 0x7f, 0x1e, 0xc8, 0xd7, 0x9e, 0xcc, 0x60, 0xbb, 0x03, 0x47,
	0xfc, 0x8a, 0x9d, 0x49, 0x6b, 0x37, 0x58, 0x4a, 0x42, 0xa3, 0xf3, 0x4f, 0xe8, 0x16, 0x41, 0x1c,
	0xac, 0x4e, 0xb2, 0xd3, 0x49, 0xfc, 0x04, 0x5d, 0xf2, 0xce, 0x2e, 0xfe, 0xfe, 0x77, 0xd6, 0x68,
	0x07, 0xfc, 0x91, 0xcd, 0x47, 0x83, 0xdc, 0x2b, 0x2c, 0x82, 0x38, 0x5c, 0x45, 0xeb, 0xa5, 0xf8,
	0xef, 0x20, 0x26, 0x0b, 0xb2, 0x99, 0x9a, 0x6c, 0x4b, 0x88, 0x45, 0x93, 0x47, 0x7e, 0xc9, 0xe6,
	0xd6, 0xb4, 0x24, 0x8b, 0x0d, 0xe4, 0x5a, 0x2a, 0x18, 0x9d, 0x66, 0xfb, 0xf0, 0x45, 0x2a, 0xe0,
	0x4b, 0x16, 0x95, 0x46, 0xd7, 0xd8, 0x0c, 0xc8, 0x91, 0x47, 0xd8, 0x10, 0xed, 0x01, 0x05, 0x15,
	0xca, 0x51, 0x2c, 0x8c, 0xc3, 0x1e, 0xf0, 0x91, 0xbf, 0xfa, 0x70, 0xf7, 0xb6, 0x6e, 0x90, 0x3e,
	0x76, 0x85, 0x28, 0x8d, 0x4a, 0xbd, 0x30, 0xe1, 0x36, 0x45, 0x5d, 0xe3, 0x77, 0x7a, 0xa8, 0xf0,
	0xe2, 0xd8, 0x97, 0x7c, 0xfb, 0x13, 0x00, 0x00, 0xff, 0xff, 0x3f, 0xb7, 0x04, 0xaa, 0x93, 0x01,
	0x00, 0x00,
}
