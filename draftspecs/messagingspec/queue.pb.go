// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: draftspecs/messagingspec/queue.proto

package messagingspec

import (
	context "context"
	envelopespec "github.com/dogmatiq/infix/draftspecs/envelopespec"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type EnqueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ApplicationKey is the identity key of the application that handles the
	// command.
	ApplicationKey string `protobuf:"bytes,1,opt,name=application_key,json=applicationKey,proto3" json:"application_key,omitempty"`
	// Envelope is the envelope containing the command to be executed.
	Envelope *envelopespec.Envelope `protobuf:"bytes,2,opt,name=envelope,proto3" json:"envelope,omitempty"`
}

func (x *EnqueueRequest) Reset() {
	*x = EnqueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueRequest) ProtoMessage() {}

func (x *EnqueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueRequest.ProtoReflect.Descriptor instead.
func (*EnqueueRequest) Descriptor() ([]byte, []int) {
	return file_draftspecs_messagingspec_queue_proto_rawDescGZIP(), []int{0}
}

func (x *EnqueueRequest) GetApplicationKey() string {
	if x != nil {
		return x.ApplicationKey
	}
	return ""
}

func (x *EnqueueRequest) GetEnvelope() *envelopespec.Envelope {
	if x != nil {
		return x.Envelope
	}
	return nil
}

type EnqueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EnqueueResponse) Reset() {
	*x = EnqueueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueResponse) ProtoMessage() {}

func (x *EnqueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueResponse.ProtoReflect.Descriptor instead.
func (*EnqueueResponse) Descriptor() ([]byte, []int) {
	return file_draftspecs_messagingspec_queue_proto_rawDescGZIP(), []int{1}
}

type AckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ApplicationKey is the identity key of the application that produced the
	// command.
	ApplicationKey string `protobuf:"bytes,1,opt,name=application_key,json=applicationKey,proto3" json:"application_key,omitempty"`
	// MessageID is the command message's unique identifier.
	MessageId string `protobuf:"bytes,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *AckRequest) Reset() {
	*x = AckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckRequest) ProtoMessage() {}

func (x *AckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckRequest.ProtoReflect.Descriptor instead.
func (*AckRequest) Descriptor() ([]byte, []int) {
	return file_draftspecs_messagingspec_queue_proto_rawDescGZIP(), []int{2}
}

func (x *AckRequest) GetApplicationKey() string {
	if x != nil {
		return x.ApplicationKey
	}
	return ""
}

func (x *AckRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type AckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AckResponse) Reset() {
	*x = AckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckResponse) ProtoMessage() {}

func (x *AckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_draftspecs_messagingspec_queue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckResponse.ProtoReflect.Descriptor instead.
func (*AckResponse) Descriptor() ([]byte, []int) {
	return file_draftspecs_messagingspec_queue_proto_rawDescGZIP(), []int{3}
}

var File_draftspecs_messagingspec_queue_proto protoreflect.FileDescriptor

var file_draftspecs_messagingspec_queue_proto_rawDesc = []byte{
	0x0a, 0x24, 0x64, 0x72, 0x61, 0x66, 0x74, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x26, 0x64, 0x72, 0x61, 0x66,
	0x74, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73,
	0x70, 0x65, 0x63, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2b, 0x64, 0x72, 0x61, 0x66, 0x74, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x72, 0x0a, 0x0e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x37, 0x0a, 0x08, 0x65, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x52, 0x08, 0x65, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x54, 0x0a, 0x0a, 0x41, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x1d, 0x0a,
	0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x0d, 0x0a, 0x0b,
	0x41, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8d, 0x02, 0x0a, 0x0c,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x52, 0x0a, 0x07,
	0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x22, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x64, 0x6f,
	0x67, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x1e, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x27, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x28, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x34, 0x5a, 0x32, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74,
	0x69, 0x71, 0x2f, 0x69, 0x6e, 0x66, 0x69, 0x78, 0x2f, 0x64, 0x72, 0x61, 0x66, 0x74, 0x73, 0x70,
	0x65, 0x63, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_draftspecs_messagingspec_queue_proto_rawDescOnce sync.Once
	file_draftspecs_messagingspec_queue_proto_rawDescData = file_draftspecs_messagingspec_queue_proto_rawDesc
)

func file_draftspecs_messagingspec_queue_proto_rawDescGZIP() []byte {
	file_draftspecs_messagingspec_queue_proto_rawDescOnce.Do(func() {
		file_draftspecs_messagingspec_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_draftspecs_messagingspec_queue_proto_rawDescData)
	})
	return file_draftspecs_messagingspec_queue_proto_rawDescData
}

var file_draftspecs_messagingspec_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_draftspecs_messagingspec_queue_proto_goTypes = []interface{}{
	(*EnqueueRequest)(nil),        // 0: dogma.messaging.v1.EnqueueRequest
	(*EnqueueResponse)(nil),       // 1: dogma.messaging.v1.EnqueueResponse
	(*AckRequest)(nil),            // 2: dogma.messaging.v1.AckRequest
	(*AckResponse)(nil),           // 3: dogma.messaging.v1.AckResponse
	(*envelopespec.Envelope)(nil), // 4: dogma.envelope.v1.Envelope
	(*MessageTypesRequest)(nil),   // 5: dogma.messaging.v1.MessageTypesRequest
	(*MessageTypesResponse)(nil),  // 6: dogma.messaging.v1.MessageTypesResponse
}
var file_draftspecs_messagingspec_queue_proto_depIdxs = []int32{
	4, // 0: dogma.messaging.v1.EnqueueRequest.envelope:type_name -> dogma.envelope.v1.Envelope
	0, // 1: dogma.messaging.v1.CommandQueue.Enqueue:input_type -> dogma.messaging.v1.EnqueueRequest
	2, // 2: dogma.messaging.v1.CommandQueue.Ack:input_type -> dogma.messaging.v1.AckRequest
	5, // 3: dogma.messaging.v1.CommandQueue.MessageTypes:input_type -> dogma.messaging.v1.MessageTypesRequest
	1, // 4: dogma.messaging.v1.CommandQueue.Enqueue:output_type -> dogma.messaging.v1.EnqueueResponse
	3, // 5: dogma.messaging.v1.CommandQueue.Ack:output_type -> dogma.messaging.v1.AckResponse
	6, // 6: dogma.messaging.v1.CommandQueue.MessageTypes:output_type -> dogma.messaging.v1.MessageTypesResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_draftspecs_messagingspec_queue_proto_init() }
func file_draftspecs_messagingspec_queue_proto_init() {
	if File_draftspecs_messagingspec_queue_proto != nil {
		return
	}
	file_draftspecs_messagingspec_messagetypes_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_draftspecs_messagingspec_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_draftspecs_messagingspec_queue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_draftspecs_messagingspec_queue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_draftspecs_messagingspec_queue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_draftspecs_messagingspec_queue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_draftspecs_messagingspec_queue_proto_goTypes,
		DependencyIndexes: file_draftspecs_messagingspec_queue_proto_depIdxs,
		MessageInfos:      file_draftspecs_messagingspec_queue_proto_msgTypes,
	}.Build()
	File_draftspecs_messagingspec_queue_proto = out.File
	file_draftspecs_messagingspec_queue_proto_rawDesc = nil
	file_draftspecs_messagingspec_queue_proto_goTypes = nil
	file_draftspecs_messagingspec_queue_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CommandQueueClient is the client API for CommandQueue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommandQueueClient interface {
	// Enqueue adds a message to the inbound command queue.
	//
	// The client is an engine with a message on its outbound queue. The server is
	// an engine that hosts an application that can handle commands of that type.
	//
	// If the server does not host the application specified in the request, it
	// MUST return a NOT_FOUND error with an attached UnrecognizedApplication
	// value.
	//
	// If the server does not handle messages of the type provided in the request,
	// it MUST return an INVALID_ARGUMENT error with an attached
	// UnrecognizedMessage value.
	//
	// The server MUST keep the enqueued message on its inbound queue until it
	// removes the message from the client's outbound queue by calling Ack(). The
	// mechanism for determining how to dial back to the client is engine-specific
	// and outside the scope of this specification.
	//
	// If the message in the request, identified by the message ID, is not already
	// on the queue, it is added. A successful response MUST be returned
	// regardless of whether the message was already enqueued or not.
	//
	// The client SHOULD retry the Enqueue() oepration until it receives a
	// successful response, or until it receives an Ack() call for the message.
	Enqueue(ctx context.Context, in *EnqueueRequest, opts ...grpc.CallOption) (*EnqueueResponse, error)
	// Ack removes a message from the outbound command queue.
	//
	// The client is an engine that has enqueued (and potentially already
	// executed) a message received via a prior call to Enqueue().
	//
	// If the message specified in the request is still on the outbound queue, it
	// is removed. A successful response MUST be returned regardless of whether
	// the message was still enqueued or not.
	//
	// Upon returning a successful response, the server MUST NOT make any future
	// call to Enqueue() for the specified message.
	Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error)
	// MessageTypes queries the messages types that the server supports for a
	// specific application.
	//
	// If the server does not host the application specified in the request, it
	// MUST return a NOT_FOUND error with an attached UnrecognizedApplication
	// value.
	MessageTypes(ctx context.Context, in *MessageTypesRequest, opts ...grpc.CallOption) (*MessageTypesResponse, error)
}

type commandQueueClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandQueueClient(cc grpc.ClientConnInterface) CommandQueueClient {
	return &commandQueueClient{cc}
}

func (c *commandQueueClient) Enqueue(ctx context.Context, in *EnqueueRequest, opts ...grpc.CallOption) (*EnqueueResponse, error) {
	out := new(EnqueueResponse)
	err := c.cc.Invoke(ctx, "/dogma.messaging.v1.CommandQueue/Enqueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandQueueClient) Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error) {
	out := new(AckResponse)
	err := c.cc.Invoke(ctx, "/dogma.messaging.v1.CommandQueue/Ack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandQueueClient) MessageTypes(ctx context.Context, in *MessageTypesRequest, opts ...grpc.CallOption) (*MessageTypesResponse, error) {
	out := new(MessageTypesResponse)
	err := c.cc.Invoke(ctx, "/dogma.messaging.v1.CommandQueue/MessageTypes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandQueueServer is the server API for CommandQueue service.
type CommandQueueServer interface {
	// Enqueue adds a message to the inbound command queue.
	//
	// The client is an engine with a message on its outbound queue. The server is
	// an engine that hosts an application that can handle commands of that type.
	//
	// If the server does not host the application specified in the request, it
	// MUST return a NOT_FOUND error with an attached UnrecognizedApplication
	// value.
	//
	// If the server does not handle messages of the type provided in the request,
	// it MUST return an INVALID_ARGUMENT error with an attached
	// UnrecognizedMessage value.
	//
	// The server MUST keep the enqueued message on its inbound queue until it
	// removes the message from the client's outbound queue by calling Ack(). The
	// mechanism for determining how to dial back to the client is engine-specific
	// and outside the scope of this specification.
	//
	// If the message in the request, identified by the message ID, is not already
	// on the queue, it is added. A successful response MUST be returned
	// regardless of whether the message was already enqueued or not.
	//
	// The client SHOULD retry the Enqueue() oepration until it receives a
	// successful response, or until it receives an Ack() call for the message.
	Enqueue(context.Context, *EnqueueRequest) (*EnqueueResponse, error)
	// Ack removes a message from the outbound command queue.
	//
	// The client is an engine that has enqueued (and potentially already
	// executed) a message received via a prior call to Enqueue().
	//
	// If the message specified in the request is still on the outbound queue, it
	// is removed. A successful response MUST be returned regardless of whether
	// the message was still enqueued or not.
	//
	// Upon returning a successful response, the server MUST NOT make any future
	// call to Enqueue() for the specified message.
	Ack(context.Context, *AckRequest) (*AckResponse, error)
	// MessageTypes queries the messages types that the server supports for a
	// specific application.
	//
	// If the server does not host the application specified in the request, it
	// MUST return a NOT_FOUND error with an attached UnrecognizedApplication
	// value.
	MessageTypes(context.Context, *MessageTypesRequest) (*MessageTypesResponse, error)
}

// UnimplementedCommandQueueServer can be embedded to have forward compatible implementations.
type UnimplementedCommandQueueServer struct {
}

func (*UnimplementedCommandQueueServer) Enqueue(context.Context, *EnqueueRequest) (*EnqueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enqueue not implemented")
}
func (*UnimplementedCommandQueueServer) Ack(context.Context, *AckRequest) (*AckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ack not implemented")
}
func (*UnimplementedCommandQueueServer) MessageTypes(context.Context, *MessageTypesRequest) (*MessageTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageTypes not implemented")
}

func RegisterCommandQueueServer(s *grpc.Server, srv CommandQueueServer) {
	s.RegisterService(&_CommandQueue_serviceDesc, srv)
}

func _CommandQueue_Enqueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnqueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandQueueServer).Enqueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dogma.messaging.v1.CommandQueue/Enqueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandQueueServer).Enqueue(ctx, req.(*EnqueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommandQueue_Ack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandQueueServer).Ack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dogma.messaging.v1.CommandQueue/Ack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandQueueServer).Ack(ctx, req.(*AckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommandQueue_MessageTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandQueueServer).MessageTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dogma.messaging.v1.CommandQueue/MessageTypes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandQueueServer).MessageTypes(ctx, req.(*MessageTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CommandQueue_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dogma.messaging.v1.CommandQueue",
	HandlerType: (*CommandQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enqueue",
			Handler:    _CommandQueue_Enqueue_Handler,
		},
		{
			MethodName: "Ack",
			Handler:    _CommandQueue_Ack_Handler,
		},
		{
			MethodName: "MessageTypes",
			Handler:    _CommandQueue_MessageTypes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "draftspecs/messagingspec/queue.proto",
}
