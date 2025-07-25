// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: telemetry.proto

package proto

import (
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

// TelemetryData represents cluster-level telemetry information
type TelemetryData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique cluster identifier (generated in-memory)
	ClusterId string `protobuf:"bytes,1,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	// SeaweedFS version
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Operating system (e.g., "linux/amd64")
	Os string `protobuf:"bytes,3,opt,name=os,proto3" json:"os,omitempty"`
	// Number of volume servers in the cluster
	VolumeServerCount int32 `protobuf:"varint,6,opt,name=volume_server_count,json=volumeServerCount,proto3" json:"volume_server_count,omitempty"`
	// Total disk usage across all volume servers (in bytes)
	TotalDiskBytes uint64 `protobuf:"varint,7,opt,name=total_disk_bytes,json=totalDiskBytes,proto3" json:"total_disk_bytes,omitempty"`
	// Total number of volumes in the cluster
	TotalVolumeCount int32 `protobuf:"varint,8,opt,name=total_volume_count,json=totalVolumeCount,proto3" json:"total_volume_count,omitempty"`
	// Number of filer servers in the cluster
	FilerCount int32 `protobuf:"varint,9,opt,name=filer_count,json=filerCount,proto3" json:"filer_count,omitempty"`
	// Number of broker servers in the cluster
	BrokerCount int32 `protobuf:"varint,10,opt,name=broker_count,json=brokerCount,proto3" json:"broker_count,omitempty"`
	// Unix timestamp when the data was collected
	Timestamp int64 `protobuf:"varint,11,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *TelemetryData) Reset() {
	*x = TelemetryData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryData) ProtoMessage() {}

func (x *TelemetryData) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryData.ProtoReflect.Descriptor instead.
func (*TelemetryData) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{0}
}

func (x *TelemetryData) GetClusterId() string {
	if x != nil {
		return x.ClusterId
	}
	return ""
}

func (x *TelemetryData) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *TelemetryData) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *TelemetryData) GetVolumeServerCount() int32 {
	if x != nil {
		return x.VolumeServerCount
	}
	return 0
}

func (x *TelemetryData) GetTotalDiskBytes() uint64 {
	if x != nil {
		return x.TotalDiskBytes
	}
	return 0
}

func (x *TelemetryData) GetTotalVolumeCount() int32 {
	if x != nil {
		return x.TotalVolumeCount
	}
	return 0
}

func (x *TelemetryData) GetFilerCount() int32 {
	if x != nil {
		return x.FilerCount
	}
	return 0
}

func (x *TelemetryData) GetBrokerCount() int32 {
	if x != nil {
		return x.BrokerCount
	}
	return 0
}

func (x *TelemetryData) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// TelemetryRequest is sent from SeaweedFS clusters to the telemetry server
type TelemetryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *TelemetryData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *TelemetryRequest) Reset() {
	*x = TelemetryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryRequest) ProtoMessage() {}

func (x *TelemetryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryRequest.ProtoReflect.Descriptor instead.
func (*TelemetryRequest) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{1}
}

func (x *TelemetryRequest) GetData() *TelemetryData {
	if x != nil {
		return x.Data
	}
	return nil
}

// TelemetryResponse is returned by the telemetry server
type TelemetryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TelemetryResponse) Reset() {
	*x = TelemetryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryResponse) ProtoMessage() {}

func (x *TelemetryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryResponse.ProtoReflect.Descriptor instead.
func (*TelemetryResponse) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{2}
}

func (x *TelemetryResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *TelemetryResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_telemetry_proto protoreflect.FileDescriptor

var file_telemetry_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x22, 0xce, 0x02, 0x0a,
	0x0d, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x76, 0x6f, 0x6c, 0x75, 0x6d,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x6b, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x6f, 0x6c, 0x75, 0x6d,
	0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x05, 0x4a, 0x04, 0x08, 0x05, 0x10, 0x06, 0x22, 0x40, 0x0a,
	0x10, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x54, 0x65, 0x6c, 0x65,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x47, 0x0a, 0x11, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x61, 0x77, 0x65, 0x65, 0x64, 0x66, 0x73,
	0x2f, 0x73, 0x65, 0x61, 0x77, 0x65, 0x65, 0x64, 0x66, 0x73, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_telemetry_proto_rawDescOnce sync.Once
	file_telemetry_proto_rawDescData = file_telemetry_proto_rawDesc
)

func file_telemetry_proto_rawDescGZIP() []byte {
	file_telemetry_proto_rawDescOnce.Do(func() {
		file_telemetry_proto_rawDescData = protoimpl.X.CompressGZIP(file_telemetry_proto_rawDescData)
	})
	return file_telemetry_proto_rawDescData
}

var file_telemetry_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_telemetry_proto_goTypes = []any{
	(*TelemetryData)(nil),     // 0: telemetry.TelemetryData
	(*TelemetryRequest)(nil),  // 1: telemetry.TelemetryRequest
	(*TelemetryResponse)(nil), // 2: telemetry.TelemetryResponse
}
var file_telemetry_proto_depIdxs = []int32{
	0, // 0: telemetry.TelemetryRequest.data:type_name -> telemetry.TelemetryData
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_telemetry_proto_init() }
func file_telemetry_proto_init() {
	if File_telemetry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_telemetry_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TelemetryData); i {
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
		file_telemetry_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TelemetryRequest); i {
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
		file_telemetry_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TelemetryResponse); i {
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
			RawDescriptor: file_telemetry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_telemetry_proto_goTypes,
		DependencyIndexes: file_telemetry_proto_depIdxs,
		MessageInfos:      file_telemetry_proto_msgTypes,
	}.Build()
	File_telemetry_proto = out.File
	file_telemetry_proto_rawDesc = nil
	file_telemetry_proto_goTypes = nil
	file_telemetry_proto_depIdxs = nil
}
