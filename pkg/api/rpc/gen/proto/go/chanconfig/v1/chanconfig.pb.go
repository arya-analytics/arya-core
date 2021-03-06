// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: chanconfig/v1/chanconfig.proto

package chanconfigv1

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

type ChannelConfig_DataType int32

const (
	ChannelConfig_FLOAT64 ChannelConfig_DataType = 0
	ChannelConfig_FLOAT32 ChannelConfig_DataType = 1
)

// Enum value maps for ChannelConfig_DataType.
var (
	ChannelConfig_DataType_name = map[int32]string{
		0: "FLOAT64",
		1: "FLOAT32",
	}
	ChannelConfig_DataType_value = map[string]int32{
		"FLOAT64": 0,
		"FLOAT32": 1,
	}
)

func (x ChannelConfig_DataType) Enum() *ChannelConfig_DataType {
	p := new(ChannelConfig_DataType)
	*p = x
	return p
}

func (x ChannelConfig_DataType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChannelConfig_DataType) Descriptor() protoreflect.EnumDescriptor {
	return file_chanconfig_v1_chanconfig_proto_enumTypes[0].Descriptor()
}

func (ChannelConfig_DataType) Type() protoreflect.EnumType {
	return &file_chanconfig_v1_chanconfig_proto_enumTypes[0]
}

func (x ChannelConfig_DataType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChannelConfig_DataType.Descriptor instead.
func (ChannelConfig_DataType) EnumDescriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{2, 0}
}

type ChannelConfig_ConflictPolicy int32

const (
	ChannelConfig_ERROR     ChannelConfig_ConflictPolicy = 0
	ChannelConfig_DISCARD   ChannelConfig_ConflictPolicy = 1
	ChannelConfig_OVERWRITE ChannelConfig_ConflictPolicy = 2
)

// Enum value maps for ChannelConfig_ConflictPolicy.
var (
	ChannelConfig_ConflictPolicy_name = map[int32]string{
		0: "ERROR",
		1: "DISCARD",
		2: "OVERWRITE",
	}
	ChannelConfig_ConflictPolicy_value = map[string]int32{
		"ERROR":     0,
		"DISCARD":   1,
		"OVERWRITE": 2,
	}
)

func (x ChannelConfig_ConflictPolicy) Enum() *ChannelConfig_ConflictPolicy {
	p := new(ChannelConfig_ConflictPolicy)
	*p = x
	return p
}

func (x ChannelConfig_ConflictPolicy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChannelConfig_ConflictPolicy) Descriptor() protoreflect.EnumDescriptor {
	return file_chanconfig_v1_chanconfig_proto_enumTypes[1].Descriptor()
}

func (ChannelConfig_ConflictPolicy) Type() protoreflect.EnumType {
	return &file_chanconfig_v1_chanconfig_proto_enumTypes[1]
}

func (x ChannelConfig_ConflictPolicy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChannelConfig_ConflictPolicy.Descriptor instead.
func (ChannelConfig_ConflictPolicy) EnumDescriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{2, 1}
}

type RetrieveConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId int32 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Limit  int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *RetrieveConfigRequest) Reset() {
	*x = RetrieveConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveConfigRequest) ProtoMessage() {}

func (x *RetrieveConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveConfigRequest.ProtoReflect.Descriptor instead.
func (*RetrieveConfigRequest) Descriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{0}
}

func (x *RetrieveConfigRequest) GetNodeId() int32 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *RetrieveConfigRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type RetrieveConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configs []*ChannelConfig `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
}

func (x *RetrieveConfigResponse) Reset() {
	*x = RetrieveConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveConfigResponse) ProtoMessage() {}

func (x *RetrieveConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveConfigResponse.ProtoReflect.Descriptor instead.
func (*RetrieveConfigResponse) Descriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{1}
}

func (x *RetrieveConfigResponse) GetConfigs() []*ChannelConfig {
	if x != nil {
		return x.Configs
	}
	return nil
}

type ChannelConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             string                       `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name           string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	NodeId         int32                        `protobuf:"varint,3,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	DataRate       float64                      `protobuf:"fixed64,4,opt,name=data_rate,json=dataRate,proto3" json:"data_rate,omitempty"`
	DataType       ChannelConfig_DataType       `protobuf:"varint,5,opt,name=data_type,json=dataType,proto3,enum=chanconfig.v1.ChannelConfig_DataType" json:"data_type,omitempty"`
	ConflictPolicy ChannelConfig_ConflictPolicy `protobuf:"varint,6,opt,name=conflict_policy,json=conflictPolicy,proto3,enum=chanconfig.v1.ChannelConfig_ConflictPolicy" json:"conflict_policy,omitempty"`
}

func (x *ChannelConfig) Reset() {
	*x = ChannelConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelConfig) ProtoMessage() {}

func (x *ChannelConfig) ProtoReflect() protoreflect.Message {
	mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelConfig.ProtoReflect.Descriptor instead.
func (*ChannelConfig) Descriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{2}
}

func (x *ChannelConfig) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *ChannelConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChannelConfig) GetNodeId() int32 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *ChannelConfig) GetDataRate() float64 {
	if x != nil {
		return x.DataRate
	}
	return 0
}

func (x *ChannelConfig) GetDataType() ChannelConfig_DataType {
	if x != nil {
		return x.DataType
	}
	return ChannelConfig_FLOAT64
}

func (x *ChannelConfig) GetConflictPolicy() ChannelConfig_ConflictPolicy {
	if x != nil {
		return x.ConflictPolicy
	}
	return ChannelConfig_ERROR
}

type CreateConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *ChannelConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *CreateConfigRequest) Reset() {
	*x = CreateConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateConfigRequest) ProtoMessage() {}

func (x *CreateConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateConfigRequest.ProtoReflect.Descriptor instead.
func (*CreateConfigRequest) Descriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{3}
}

func (x *CreateConfigRequest) GetConfig() *ChannelConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

type CreateConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateConfigResponse) Reset() {
	*x = CreateConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateConfigResponse) ProtoMessage() {}

func (x *CreateConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chanconfig_v1_chanconfig_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateConfigResponse.ProtoReflect.Descriptor instead.
func (*CreateConfigResponse) Descriptor() ([]byte, []int) {
	return file_chanconfig_v1_chanconfig_proto_rawDescGZIP(), []int{4}
}

var File_chanconfig_v1_chanconfig_proto protoreflect.FileDescriptor

var file_chanconfig_v1_chanconfig_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x22,
	0x46, 0x0a, 0x15, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x50, 0x0a, 0x16, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x36, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0xe2, 0x02, 0x0a, 0x0d, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x61, 0x74, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x54, 0x0a, 0x0f, 0x63, 0x6f, 0x6e,
	0x66, 0x6c, 0x69, 0x63, 0x74, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52,
	0x0e, 0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22,
	0x24, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x46,
	0x4c, 0x4f, 0x41, 0x54, 0x36, 0x34, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x46, 0x4c, 0x4f, 0x41,
	0x54, 0x33, 0x32, 0x10, 0x01, 0x22, 0x37, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63,
	0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x49, 0x53, 0x43, 0x41, 0x52, 0x44, 0x10, 0x01, 0x12,
	0x0d, 0x0a, 0x09, 0x4f, 0x56, 0x45, 0x52, 0x57, 0x52, 0x49, 0x54, 0x45, 0x10, 0x02, 0x22, 0x4b,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x16, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xcb, 0x01, 0x0a, 0x11, 0x43, 0x68, 0x61, 0x6e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x22, 0x2e, 0x63, 0x68, 0x61, 0x6e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x5d, 0x0a, 0x0e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x24, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x63, 0x68, 0x61,
	0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0xd2, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x43, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x57, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x79, 0x61, 0x2d, 0x61, 0x6e, 0x61, 0x6c,
	0x79, 0x74, 0x69, 0x63, 0x73, 0x2f, 0x61, 0x72, 0x79, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x0d, 0x43, 0x68, 0x61, 0x6e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0d, 0x43, 0x68, 0x61, 0x6e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x19, 0x43, 0x68, 0x61, 0x6e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chanconfig_v1_chanconfig_proto_rawDescOnce sync.Once
	file_chanconfig_v1_chanconfig_proto_rawDescData = file_chanconfig_v1_chanconfig_proto_rawDesc
)

func file_chanconfig_v1_chanconfig_proto_rawDescGZIP() []byte {
	file_chanconfig_v1_chanconfig_proto_rawDescOnce.Do(func() {
		file_chanconfig_v1_chanconfig_proto_rawDescData = protoimpl.X.CompressGZIP(file_chanconfig_v1_chanconfig_proto_rawDescData)
	})
	return file_chanconfig_v1_chanconfig_proto_rawDescData
}

var file_chanconfig_v1_chanconfig_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_chanconfig_v1_chanconfig_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_chanconfig_v1_chanconfig_proto_goTypes = []interface{}{
	(ChannelConfig_DataType)(0),       // 0: chanconfig.v1.ChannelConfig.DataType
	(ChannelConfig_ConflictPolicy)(0), // 1: chanconfig.v1.ChannelConfig.ConflictPolicy
	(*RetrieveConfigRequest)(nil),     // 2: chanconfig.v1.RetrieveConfigRequest
	(*RetrieveConfigResponse)(nil),    // 3: chanconfig.v1.RetrieveConfigResponse
	(*ChannelConfig)(nil),             // 4: chanconfig.v1.ChannelConfig
	(*CreateConfigRequest)(nil),       // 5: chanconfig.v1.CreateConfigRequest
	(*CreateConfigResponse)(nil),      // 6: chanconfig.v1.CreateConfigResponse
}
var file_chanconfig_v1_chanconfig_proto_depIdxs = []int32{
	4, // 0: chanconfig.v1.RetrieveConfigResponse.configs:type_name -> chanconfig.v1.ChannelConfig
	0, // 1: chanconfig.v1.ChannelConfig.data_type:type_name -> chanconfig.v1.ChannelConfig.DataType
	1, // 2: chanconfig.v1.ChannelConfig.conflict_policy:type_name -> chanconfig.v1.ChannelConfig.ConflictPolicy
	4, // 3: chanconfig.v1.CreateConfigRequest.config:type_name -> chanconfig.v1.ChannelConfig
	5, // 4: chanconfig.v1.ChanConfigService.CreateConfig:input_type -> chanconfig.v1.CreateConfigRequest
	2, // 5: chanconfig.v1.ChanConfigService.RetrieveConfig:input_type -> chanconfig.v1.RetrieveConfigRequest
	6, // 6: chanconfig.v1.ChanConfigService.CreateConfig:output_type -> chanconfig.v1.CreateConfigResponse
	3, // 7: chanconfig.v1.ChanConfigService.RetrieveConfig:output_type -> chanconfig.v1.RetrieveConfigResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_chanconfig_v1_chanconfig_proto_init() }
func file_chanconfig_v1_chanconfig_proto_init() {
	if File_chanconfig_v1_chanconfig_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chanconfig_v1_chanconfig_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveConfigRequest); i {
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
		file_chanconfig_v1_chanconfig_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveConfigResponse); i {
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
		file_chanconfig_v1_chanconfig_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelConfig); i {
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
		file_chanconfig_v1_chanconfig_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateConfigRequest); i {
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
		file_chanconfig_v1_chanconfig_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateConfigResponse); i {
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
			RawDescriptor: file_chanconfig_v1_chanconfig_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chanconfig_v1_chanconfig_proto_goTypes,
		DependencyIndexes: file_chanconfig_v1_chanconfig_proto_depIdxs,
		EnumInfos:         file_chanconfig_v1_chanconfig_proto_enumTypes,
		MessageInfos:      file_chanconfig_v1_chanconfig_proto_msgTypes,
	}.Build()
	File_chanconfig_v1_chanconfig_proto = out.File
	file_chanconfig_v1_chanconfig_proto_rawDesc = nil
	file_chanconfig_v1_chanconfig_proto_goTypes = nil
	file_chanconfig_v1_chanconfig_proto_depIdxs = nil
}
