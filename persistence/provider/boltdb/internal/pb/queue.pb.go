// Code generated by protoc-gen-go. DO NOT EDIT.
// source: persistence/provider/boltdb/internal/pb/queue.proto

package pb

import (
	fmt "fmt"
	envelopespec "github.com/dogmatiq/infix/draftspecs/envelopespec"
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

// QueueMessage is a protocol buffers representation of an Infix queue.Message.
type QueueMessage struct {
	Revision             uint64                 `protobuf:"varint,1,opt,name=revision,proto3" json:"revision,omitempty"`
	Envelope             *envelopespec.Envelope `protobuf:"bytes,2,opt,name=Envelope,proto3" json:"Envelope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *QueueMessage) Reset()         { *m = QueueMessage{} }
func (m *QueueMessage) String() string { return proto.CompactTextString(m) }
func (*QueueMessage) ProtoMessage()    {}
func (*QueueMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_910269091ef3cb55, []int{0}
}

func (m *QueueMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueueMessage.Unmarshal(m, b)
}
func (m *QueueMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueueMessage.Marshal(b, m, deterministic)
}
func (m *QueueMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueueMessage.Merge(m, src)
}
func (m *QueueMessage) XXX_Size() int {
	return xxx_messageInfo_QueueMessage.Size(m)
}
func (m *QueueMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_QueueMessage.DiscardUnknown(m)
}

var xxx_messageInfo_QueueMessage proto.InternalMessageInfo

func (m *QueueMessage) GetRevision() uint64 {
	if m != nil {
		return m.Revision
	}
	return 0
}

func (m *QueueMessage) GetEnvelope() *envelopespec.Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

func init() {
	proto.RegisterType((*QueueMessage)(nil), "infix.persistence.boltdb.v1.QueueMessage")
}

func init() {
	proto.RegisterFile("persistence/provider/boltdb/internal/pb/queue.proto", fileDescriptor_910269091ef3cb55)
}

var fileDescriptor_910269091ef3cb55 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x8f, 0x3f, 0x4b, 0x03, 0x41,
	0x10, 0x47, 0x39, 0x11, 0x09, 0xab, 0xd5, 0x55, 0x21, 0x69, 0x82, 0x85, 0xa4, 0x9a, 0x21, 0xa6,
	0xb0, 0x56, 0xb1, 0xb4, 0x30, 0xa5, 0xdd, 0xfe, 0x99, 0x9c, 0x03, 0x97, 0x9d, 0xcd, 0xee, 0xde,
	0xe2, 0xc7, 0x97, 0xbb, 0xd3, 0xc5, 0x32, 0xe5, 0x83, 0x79, 0x8f, 0xdf, 0xa8, 0x7d, 0xa0, 0x98,
	0x38, 0x65, 0xf2, 0x96, 0x30, 0x44, 0x29, 0xec, 0x28, 0xa2, 0x91, 0x3e, 0x3b, 0x83, 0xec, 0x33,
	0x45, 0xaf, 0x7b, 0x0c, 0x06, 0xcf, 0x03, 0x0d, 0x04, 0x21, 0x4a, 0x96, 0x76, 0xcd, 0xfe, 0xc8,
	0xdf, 0xf0, 0x4f, 0x85, 0xd9, 0x80, 0xb2, 0x5b, 0x3d, 0xb8, 0xa8, 0x8f, 0x39, 0x05, 0xb2, 0x09,
	0xc9, 0x17, 0xea, 0x25, 0xd0, 0x48, 0x15, 0xe6, 0xc8, 0xbd, 0x55, 0x77, 0x1f, 0x63, 0xf3, 0x9d,
	0x52, 0xd2, 0x1d, 0xb5, 0x2b, 0xb5, 0x88, 0x54, 0x38, 0xb1, 0xf8, 0x65, 0xb3, 0x69, 0xb6, 0xd7,
	0x87, 0xca, 0xed, 0x93, 0x5a, 0xbc, 0xfd, 0xda, 0xcb, 0xab, 0x4d, 0xb3, 0xbd, 0x7d, 0x5c, 0x83,
	0x93, 0xee, 0xa4, 0xa1, 0x46, 0xcb, 0x0e, 0xfe, 0x4e, 0x0e, 0xf5, 0xf8, 0xe5, 0xf5, 0xf3, 0xb9,
	0xe3, 0xfc, 0x35, 0x18, 0xb0, 0x72, 0xc2, 0x49, 0xc9, 0x7c, 0xc6, 0x69, 0x3f, 0x5e, 0xf8, 0xba,
	0xb9, 0x99, 0x06, 0xef, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x53, 0x85, 0xbe, 0x2c, 0x01,
	0x00, 0x00,
}
