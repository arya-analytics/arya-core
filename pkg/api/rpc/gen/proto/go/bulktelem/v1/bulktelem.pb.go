// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: bulktelem/v1/bulktelem.proto

package bulktelemv1

import (
	rpc "github.com/arya-analytics/arya-core/pkg/api/rpc/gen/proto/go/google/rpc"
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

type CreateStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelConfigId string `protobuf:"bytes,1,opt,name=channel_config_id,json=channelConfigId,proto3" json:"channel_config_id,omitempty"`
	StartTs         int64  `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	Data            []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *CreateStreamRequest) Reset() {
	*x = CreateStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStreamRequest) ProtoMessage() {}

func (x *CreateStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStreamRequest.ProtoReflect.Descriptor instead.
func (*CreateStreamRequest) Descriptor() ([]byte, []int) {
	return file_bulktelem_v1_bulktelem_proto_rawDescGZIP(), []int{0}
}

func (x *CreateStreamRequest) GetChannelConfigId() string {
	if x != nil {
		return x.ChannelConfigId
	}
	return ""
}

func (x *CreateStreamRequest) GetStartTs() int64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *CreateStreamRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type CreateStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *CreateStreamResponse) Reset() {
	*x = CreateStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStreamResponse) ProtoMessage() {}

func (x *CreateStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStreamResponse.ProtoReflect.Descriptor instead.
func (*CreateStreamResponse) Descriptor() ([]byte, []int) {
	return file_bulktelem_v1_bulktelem_proto_rawDescGZIP(), []int{1}
}

func (x *CreateStreamResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       rpc.Code `protobuf:"varint,1,opt,name=type,proto3,enum=google.rpc.Code" json:"type,omitempty"`
	TypeString string   `protobuf:"bytes,2,opt,name=type_string,json=typeString,proto3" json:"type_string,omitempty"`
	Message    string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_bulktelem_v1_bulktelem_proto_rawDescGZIP(), []int{2}
}

func (x *Error) GetType() rpc.Code {
	if x != nil {
		return x.Type
	}
	return rpc.Code_OK
}

func (x *Error) GetTypeString() string {
	if x != nil {
		return x.TypeString
	}
	return ""
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type RetrieveStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelConfigId string `protobuf:"bytes,1,opt,name=channel_config_id,json=channelConfigId,proto3" json:"channel_config_id,omitempty"`
	StartTs         int64  `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	EndTs           int64  `protobuf:"varint,3,opt,name=end_ts,json=endTs,proto3" json:"end_ts,omitempty"`
}

func (x *RetrieveStreamRequest) Reset() {
	*x = RetrieveStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveStreamRequest) ProtoMessage() {}

func (x *RetrieveStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveStreamRequest.ProtoReflect.Descriptor instead.
func (*RetrieveStreamRequest) Descriptor() ([]byte, []int) {
	return file_bulktelem_v1_bulktelem_proto_rawDescGZIP(), []int{3}
}

func (x *RetrieveStreamRequest) GetChannelConfigId() string {
	if x != nil {
		return x.ChannelConfigId
	}
	return ""
}

func (x *RetrieveStreamRequest) GetStartTs() int64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *RetrieveStreamRequest) GetEndTs() int64 {
	if x != nil {
		return x.EndTs
	}
	return 0
}

type RetrieveStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTs  int64   `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	DataRate float32 `protobuf:"fixed32,3,opt,name=data_rate,json=dataRate,proto3" json:"data_rate,omitempty"`
	DataType int64   `protobuf:"varint,4,opt,name=data_type,json=dataType,proto3" json:"data_type,omitempty"`
	Data     []byte  `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RetrieveStreamResponse) Reset() {
	*x = RetrieveStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveStreamResponse) ProtoMessage() {}

func (x *RetrieveStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bulktelem_v1_bulktelem_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveStreamResponse.ProtoReflect.Descriptor instead.
func (*RetrieveStreamResponse) Descriptor() ([]byte, []int) {
	return file_bulktelem_v1_bulktelem_proto_rawDescGZIP(), []int{4}
}

func (x *RetrieveStreamResponse) GetStartTs() int64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *RetrieveStreamResponse) GetDataRate() float32 {
	if x != nil {
		return x.DataRate
	}
	return 0
}

func (x *RetrieveStreamResponse) GetDataType() int64 {
	if x != nil {
		return x.DataType
	}
	return 0
}

func (x *RetrieveStreamResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_bulktelem_v1_bulktelem_proto protoreflect.FileDescriptor

var file_bulktelem_v1_bulktelem_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x62,
	0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x15, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x41, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a,
	0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x62,
	0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x68, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x10, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x5f,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x79,
	0x70, 0x65, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x75, 0x0a, 0x15, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x54, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x16, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x52, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xcc, 0x01,
	0x0a, 0x10, 0x42, 0x75, 0x6c, 0x6b, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x59, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x21, 0x2e, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x5d, 0x0a,
	0x0e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x23, 0x2e, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0xca, 0x01, 0x0a,
	0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x76,
	0x31, 0x42, 0x0e, 0x42, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x72, 0x79, 0x61, 0x2d, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2f, 0x61,
	0x72, 0x79, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x2f, 0x62, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x62,
	0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x42, 0x58, 0x58,
	0xaa, 0x02, 0x0c, 0x42, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x0c, 0x42, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x18, 0x42, 0x75, 0x6c, 0x6b, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x42, 0x75, 0x6c, 0x6b,
	0x74, 0x65, 0x6c, 0x65, 0x6d, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_bulktelem_v1_bulktelem_proto_rawDescOnce sync.Once
	file_bulktelem_v1_bulktelem_proto_rawDescData = file_bulktelem_v1_bulktelem_proto_rawDesc
)

func file_bulktelem_v1_bulktelem_proto_rawDescGZIP() []byte {
	file_bulktelem_v1_bulktelem_proto_rawDescOnce.Do(func() {
		file_bulktelem_v1_bulktelem_proto_rawDescData = protoimpl.X.CompressGZIP(file_bulktelem_v1_bulktelem_proto_rawDescData)
	})
	return file_bulktelem_v1_bulktelem_proto_rawDescData
}

var file_bulktelem_v1_bulktelem_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_bulktelem_v1_bulktelem_proto_goTypes = []interface{}{
	(*CreateStreamRequest)(nil),    // 0: bulktelem.v1.CreateStreamRequest
	(*CreateStreamResponse)(nil),   // 1: bulktelem.v1.CreateStreamResponse
	(*Error)(nil),                  // 2: bulktelem.v1.error
	(*RetrieveStreamRequest)(nil),  // 3: bulktelem.v1.RetrieveStreamRequest
	(*RetrieveStreamResponse)(nil), // 4: bulktelem.v1.RetrieveStreamResponse
	(rpc.Code)(0),                  // 5: google.rpc.Code
}
var file_bulktelem_v1_bulktelem_proto_depIdxs = []int32{
	2, // 0: bulktelem.v1.CreateStreamResponse.Error:type_name -> bulktelem.v1.error
	5, // 1: bulktelem.v1.error.type:type_name -> google.rpc.Code
	0, // 2: bulktelem.v1.BulkTelemService.CreateStream:input_type -> bulktelem.v1.CreateStreamRequest
	3, // 3: bulktelem.v1.BulkTelemService.RetrieveStream:input_type -> bulktelem.v1.RetrieveStreamRequest
	1, // 4: bulktelem.v1.BulkTelemService.CreateStream:output_type -> bulktelem.v1.CreateStreamResponse
	4, // 5: bulktelem.v1.BulkTelemService.RetrieveStream:output_type -> bulktelem.v1.RetrieveStreamResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_bulktelem_v1_bulktelem_proto_init() }
func file_bulktelem_v1_bulktelem_proto_init() {
	if File_bulktelem_v1_bulktelem_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bulktelem_v1_bulktelem_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStreamRequest); i {
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
		file_bulktelem_v1_bulktelem_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStreamResponse); i {
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
		file_bulktelem_v1_bulktelem_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_bulktelem_v1_bulktelem_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveStreamRequest); i {
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
		file_bulktelem_v1_bulktelem_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveStreamResponse); i {
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
			RawDescriptor: file_bulktelem_v1_bulktelem_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_bulktelem_v1_bulktelem_proto_goTypes,
		DependencyIndexes: file_bulktelem_v1_bulktelem_proto_depIdxs,
		MessageInfos:      file_bulktelem_v1_bulktelem_proto_msgTypes,
	}.Build()
	File_bulktelem_v1_bulktelem_proto = out.File
	file_bulktelem_v1_bulktelem_proto_rawDesc = nil
	file_bulktelem_v1_bulktelem_proto_goTypes = nil
	file_bulktelem_v1_bulktelem_proto_depIdxs = nil
}
