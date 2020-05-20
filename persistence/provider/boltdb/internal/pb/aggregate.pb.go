// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
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

// AggregateMetaData is a protocol buffers representation of
// persistence.AggregateMetaData.
type AggregateMetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Revision        uint64 `protobuf:"varint,1,opt,name=revision,proto3" json:"revision,omitempty"`
	InstanceExists  bool   `protobuf:"varint,2,opt,name=instance_exists,json=instanceExists,proto3" json:"instance_exists,omitempty"`
	LastDestroyedBy string `protobuf:"bytes,3,opt,name=last_destroyed_by,json=lastDestroyedBy,proto3" json:"last_destroyed_by,omitempty"`
}

func (x *AggregateMetaData) Reset() {
	*x = AggregateMetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_persistence_provider_boltdb_internal_pb_aggregate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateMetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateMetaData) ProtoMessage() {}

func (x *AggregateMetaData) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use AggregateMetaData.ProtoReflect.Descriptor instead.
func (*AggregateMetaData) Descriptor() ([]byte, []int) {
	return file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDescGZIP(), []int{0}
}

func (x *AggregateMetaData) GetRevision() uint64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *AggregateMetaData) GetInstanceExists() bool {
	if x != nil {
		return x.InstanceExists
	}
	return false
}

func (x *AggregateMetaData) GetLastDestroyedBy() string {
	if x != nil {
		return x.LastDestroyedBy
	}
	return ""
}

var File_persistence_provider_boltdb_internal_pb_aggregate_proto protoreflect.FileDescriptor

var file_persistence_provider_boltdb_internal_pb_aggregate_proto_rawDesc = []byte{
	0x0a, 0x37, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f, 0x62, 0x6f, 0x6c, 0x74, 0x64, 0x62, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x69, 0x6e, 0x66, 0x69, 0x78,
	0x2e, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x62, 0x6f, 0x6c,
	0x74, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x22, 0x84, 0x01, 0x0a, 0x11, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74,
	0x73, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f,
	0x79, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6c, 0x61,
	0x73, 0x74, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x42, 0x79, 0x42, 0x43, 0x5a,
	0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x74, 0x69, 0x71, 0x2f, 0x69, 0x6e, 0x66, 0x69, 0x78, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f,
	0x62, 0x6f, 0x6c, 0x74, 0x64, 0x62, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
	(*AggregateMetaData)(nil), // 0: infix.persistence.boltdb.v1.AggregateMetaData
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
			switch v := v.(*AggregateMetaData); i {
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
