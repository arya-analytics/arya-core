// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: chanstream/v1/chanstream.proto

package chanstreamv1

import (
	v1 "github.com/arya-analytics/aryacore/pkg/rpc/gen/proto/go/error/v1"
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

type ChannelSample struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelConfigId string  `protobuf:"bytes,1,opt,name=channel_config_id,json=channelConfigId,proto3" json:"channel_config_id,omitempty"`
	Value           float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp       int64   `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ChannelSample) Reset() {
	*x = ChannelSample{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanstream_v1_chanstream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelSample) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelSample) ProtoMessage() {}

func (x *ChannelSample) ProtoReflect() protoreflect.Message {
	mi := &file_chanstream_v1_chanstream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelSample.ProtoReflect.Descriptor instead.
func (*ChannelSample) Descriptor() ([]byte, []int) {
	return file_chanstream_v1_chanstream_proto_rawDescGZIP(), []int{0}
}

func (x *ChannelSample) GetChannelConfigId() string {
	if x != nil {
		return x.ChannelConfigId
	}
	return ""
}

func (x *ChannelSample) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *ChannelSample) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelSample *ChannelSample `protobuf:"bytes,1,opt,name=channel_sample,json=channelSample,proto3" json:"channel_sample,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanstream_v1_chanstream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chanstream_v1_chanstream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_chanstream_v1_chanstream_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRequest) GetChannelSample() *ChannelSample {
	if x != nil {
		return x.ChannelSample
	}
	return nil
}

type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *v1.Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanstream_v1_chanstream_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chanstream_v1_chanstream_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_chanstream_v1_chanstream_proto_rawDescGZIP(), []int{2}
}

func (x *CreateResponse) GetError() *v1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type RetrieveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pkc []string `protobuf:"bytes,1,rep,name=pkc,proto3" json:"pkc,omitempty"`
}

func (x *RetrieveRequest) Reset() {
	*x = RetrieveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanstream_v1_chanstream_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveRequest) ProtoMessage() {}

func (x *RetrieveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chanstream_v1_chanstream_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveRequest.ProtoReflect.Descriptor instead.
func (*RetrieveRequest) Descriptor() ([]byte, []int) {
	return file_chanstream_v1_chanstream_proto_rawDescGZIP(), []int{3}
}

func (x *RetrieveRequest) GetPkc() []string {
	if x != nil {
		return x.Pkc
	}
	return nil
}

type RetrieveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelSample *ChannelSample `protobuf:"bytes,1,opt,name=channel_sample,json=channelSample,proto3" json:"channel_sample,omitempty"`
	Error         *v1.Error      `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *RetrieveResponse) Reset() {
	*x = RetrieveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanstream_v1_chanstream_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveResponse) ProtoMessage() {}

func (x *RetrieveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chanstream_v1_chanstream_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveResponse.ProtoReflect.Descriptor instead.
func (*RetrieveResponse) Descriptor() ([]byte, []int) {
	return file_chanstream_v1_chanstream_proto_rawDescGZIP(), []int{4}
}

func (x *RetrieveResponse) GetChannelSample() *ChannelSample {
	if x != nil {
		return x.ChannelSample
	}
	return nil
}

func (x *RetrieveResponse) GetError() *v1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_chanstream_v1_chanstream_proto protoreflect.FileDescriptor

var file_chanstream_v1_chanstream_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x1a,
	0x14, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0d, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x54, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x43, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x5f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x0d, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x22, 0x37, 0x0a, 0x0e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x23, 0x0a, 0x0f, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6b, 0x63, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x70, 0x6b, 0x63, 0x22, 0x7e, 0x0a, 0x10, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43,
	0x0a, 0x0e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xb2, 0x01, 0x0a, 0x14, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e,
	0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x68,
	0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x4f,
	0x0a, 0x08, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x12, 0x1e, 0x2e, 0x63, 0x68, 0x61,
	0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x68, 0x61,
	0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42,
	0xd2, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x43, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x57, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x79, 0x61, 0x2d, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74,
	0x69, 0x63, 0x73, 0x2f, 0x61, 0x72, 0x79, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x0d, 0x43, 0x68, 0x61, 0x6e, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0d, 0x43, 0x68, 0x61, 0x6e, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x19, 0x43, 0x68, 0x61, 0x6e, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chanstream_v1_chanstream_proto_rawDescOnce sync.Once
	file_chanstream_v1_chanstream_proto_rawDescData = file_chanstream_v1_chanstream_proto_rawDesc
)

func file_chanstream_v1_chanstream_proto_rawDescGZIP() []byte {
	file_chanstream_v1_chanstream_proto_rawDescOnce.Do(func() {
		file_chanstream_v1_chanstream_proto_rawDescData = protoimpl.X.CompressGZIP(file_chanstream_v1_chanstream_proto_rawDescData)
	})
	return file_chanstream_v1_chanstream_proto_rawDescData
}

var file_chanstream_v1_chanstream_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_chanstream_v1_chanstream_proto_goTypes = []interface{}{
	(*ChannelSample)(nil),    // 0: chanstream.v1.ChannelSample
	(*CreateRequest)(nil),    // 1: chanstream.v1.CreateRequest
	(*CreateResponse)(nil),   // 2: chanstream.v1.CreateResponse
	(*RetrieveRequest)(nil),  // 3: chanstream.v1.RetrieveRequest
	(*RetrieveResponse)(nil), // 4: chanstream.v1.RetrieveResponse
	(*v1.Error)(nil),         // 5: error.v1.Error
}
var file_chanstream_v1_chanstream_proto_depIdxs = []int32{
	0, // 0: chanstream.v1.CreateRequest.channel_sample:type_name -> chanstream.v1.ChannelSample
	5, // 1: chanstream.v1.CreateResponse.error:type_name -> error.v1.Error
	0, // 2: chanstream.v1.RetrieveResponse.channel_sample:type_name -> chanstream.v1.ChannelSample
	5, // 3: chanstream.v1.RetrieveResponse.error:type_name -> error.v1.Error
	1, // 4: chanstream.v1.ChannelStreamService.Create:input_type -> chanstream.v1.CreateRequest
	3, // 5: chanstream.v1.ChannelStreamService.Retrieve:input_type -> chanstream.v1.RetrieveRequest
	2, // 6: chanstream.v1.ChannelStreamService.Create:output_type -> chanstream.v1.CreateResponse
	4, // 7: chanstream.v1.ChannelStreamService.Retrieve:output_type -> chanstream.v1.RetrieveResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_chanstream_v1_chanstream_proto_init() }
func file_chanstream_v1_chanstream_proto_init() {
	if File_chanstream_v1_chanstream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chanstream_v1_chanstream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelSample); i {
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
		file_chanstream_v1_chanstream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_chanstream_v1_chanstream_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateResponse); i {
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
		file_chanstream_v1_chanstream_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveRequest); i {
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
		file_chanstream_v1_chanstream_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveResponse); i {
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
			RawDescriptor: file_chanstream_v1_chanstream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chanstream_v1_chanstream_proto_goTypes,
		DependencyIndexes: file_chanstream_v1_chanstream_proto_depIdxs,
		MessageInfos:      file_chanstream_v1_chanstream_proto_msgTypes,
	}.Build()
	File_chanstream_v1_chanstream_proto = out.File
	file_chanstream_v1_chanstream_proto_rawDesc = nil
	file_chanstream_v1_chanstream_proto_goTypes = nil
	file_chanstream_v1_chanstream_proto_depIdxs = nil
}
