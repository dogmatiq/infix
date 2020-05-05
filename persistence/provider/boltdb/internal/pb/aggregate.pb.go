// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: persistence/provider/boltdb/internal/pb/aggregate.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
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

// AggregateStoreMetaData is a protocol buffers representation of
// aggregatestore.MetaData.
type AggregateStoreMetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Revision    uint64 `protobuf:"varint,1,opt,name=revision,proto3" json:"revision,omitempty"`
	BeginOffset uint64 `protobuf:"varint,2,opt,name=begin_offset,json=beginOffset,proto3" json:"begin_offset,omitempty"`
	EndOffset   uint64 `protobuf:"varint,3,opt,name=end_offset,json=endOffset,proto3" json:"end_offset,omitempty"`
}

func (x *AggregateStoreMetaData) Reset() {
	*x = AggregateStoreMetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateStoreMetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateStoreMetaData) ProtoMessage() {}

func (x *AggregateStoreMetaData) ProtoReflect() protoreflect.Message {
	mi := &file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregateStoreMetaData.ProtoReflect.Descriptor instead.
func (*AggregateStoreMetaData) Descriptor() ([]byte, []int) {
	return file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescGZIP(), []int{0}
}

func (x *AggregateStoreMetaData) GetRevision() uint64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *AggregateStoreMetaData) GetBeginOffset() uint64 {
	if x != nil {
		return x.BeginOffset
	}
	return 0
}

func (x *AggregateStoreMetaData) GetEndOffset() uint64 {
	if x != nil {
		return x.EndOffset
	}
	return 0
}

var File_persistence_provider_boltdb_internal_pb_aggregate_proto protoreflect.FileDescriptor

var file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDesc = []byte{
	0x0a, 0x37, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f, 0x62, 0x6f, 0x6c, 0x74, 0x64, 0x62, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x69, 0x6e, 0x66, 0x69, 0x78,
	0x2e, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x62, 0x6f, 0x6c,
	0x74, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x22, 0x76, 0x0a, 0x16, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x65, 0x67, 0x69, 0x6e, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0b, 0x62, 0x65, 0x67, 0x69, 0x6e, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x65, 0x6e, 0x64, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x42, 0x43,
	0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x69, 0x6e, 0x66, 0x69, 0x78, 0x2f, 0x70, 0x65, 0x72, 0x73,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x2f, 0x62, 0x6f, 0x6c, 0x74, 0x64, 0x62, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescOnce sync.Once
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescData = file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDesc
)

func file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescGZIP() []byte {
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescOnce.Do(func() {
		file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescData = protoimpl.X.CompressGZIP(file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescData)
	})
	return file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescData
}

var file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_persistence_provider_boltdb_internal_pb_aggregate_proto_goTypes = []interface{}{
	(*AggregateStoreMetaData)(nil), // 0: infix.persistence.boltdb.v1.AggregateStoreMetaData
}
var file_persistence_provider_boltdb_internal_pb_aggregate_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_persistence_provider_boltdb_internal_pb_aggregate_proto_init() }
func file_persistence_provider_boltdb_internal_pb_aggregate_proto_init() {
	if File_persistence_provider_boltdb_internal_pb_aggregate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregateStoreMetaData); i {
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
			RawDescriptor: file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_persistence_provider_boltdb_internal_pb_aggregate_proto_goTypes,
		DependencyIndexes: file_persistence_provider_boltdb_internal_pb_aggregate_proto_depIdxs,
		MessageInfos:      file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes,
	}.Build()
	File_persistence_provider_boltdb_internal_pb_aggregate_proto = out.File
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDesc = nil
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_goTypes = nil
	file_persistence_provider_boltdb_internal_pb_aggregate_proto_depIdxs = nil
}
