// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: contact.proto

package contact

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

// 空消息
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{0}
}

// 消息结构体
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId  int64  `protobuf:"varint,1,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	Content    string `protobuf:"bytes,2,opt,name=Content,proto3" json:"Content,omitempty"`
	CreateTime int64  `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	FromId     int64  `protobuf:"varint,4,opt,name=FromId,proto3" json:"FromId,omitempty"`
	ToId       int64  `protobuf:"varint,5,opt,name=ToId,proto3" json:"ToId,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *Message) GetFromId() int64 {
	if x != nil {
		return x.FromId
	}
	return 0
}

func (x *Message) GetToId() int64 {
	if x != nil {
		return x.ToId
	}
	return 0
}

type GetLatestMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserAId int64 `protobuf:"varint,1,opt,name=UserAId,proto3" json:"UserAId,omitempty"`
	UserBId int64 `protobuf:"varint,2,opt,name=UserBId,proto3" json:"UserBId,omitempty"`
}

func (x *GetLatestMessageRequest) Reset() {
	*x = GetLatestMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLatestMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLatestMessageRequest) ProtoMessage() {}

func (x *GetLatestMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLatestMessageRequest.ProtoReflect.Descriptor instead.
func (*GetLatestMessageRequest) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{2}
}

func (x *GetLatestMessageRequest) GetUserAId() int64 {
	if x != nil {
		return x.UserAId
	}
	return 0
}

func (x *GetLatestMessageRequest) GetUserBId() int64 {
	if x != nil {
		return x.UserBId
	}
	return 0
}

type GetLatestMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *Message `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *GetLatestMessageResponse) Reset() {
	*x = GetLatestMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLatestMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLatestMessageResponse) ProtoMessage() {}

func (x *GetLatestMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLatestMessageResponse.ProtoReflect.Descriptor instead.
func (*GetLatestMessageResponse) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{3}
}

func (x *GetLatestMessageResponse) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type CreateMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromId  int64  `protobuf:"varint,1,opt,name=FromId,proto3" json:"FromId,omitempty"`
	ToId    int64  `protobuf:"varint,2,opt,name=ToId,proto3" json:"ToId,omitempty"`
	Content string `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *CreateMessageRequest) Reset() {
	*x = CreateMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMessageRequest) ProtoMessage() {}

func (x *CreateMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMessageRequest.ProtoReflect.Descriptor instead.
func (*CreateMessageRequest) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{4}
}

func (x *CreateMessageRequest) GetFromId() int64 {
	if x != nil {
		return x.FromId
	}
	return 0
}

func (x *CreateMessageRequest) GetToId() int64 {
	if x != nil {
		return x.ToId
	}
	return 0
}

func (x *CreateMessageRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type GetMessageListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromId     int64 `protobuf:"varint,1,opt,name=FromId,proto3" json:"FromId,omitempty"`
	ToId       int64 `protobuf:"varint,2,opt,name=ToId,proto3" json:"ToId,omitempty"`
	PreMsgTime int64 `protobuf:"varint,3,opt,name=PreMsgTime,proto3" json:"PreMsgTime,omitempty"`
}

func (x *GetMessageListRequest) Reset() {
	*x = GetMessageListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMessageListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageListRequest) ProtoMessage() {}

func (x *GetMessageListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageListRequest.ProtoReflect.Descriptor instead.
func (*GetMessageListRequest) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{5}
}

func (x *GetMessageListRequest) GetFromId() int64 {
	if x != nil {
		return x.FromId
	}
	return 0
}

func (x *GetMessageListRequest) GetToId() int64 {
	if x != nil {
		return x.ToId
	}
	return 0
}

func (x *GetMessageListRequest) GetPreMsgTime() int64 {
	if x != nil {
		return x.PreMsgTime
	}
	return 0
}

type GetMessageListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*Message `protobuf:"bytes,1,rep,name=Messages,proto3" json:"Messages,omitempty"`
}

func (x *GetMessageListResponse) Reset() {
	*x = GetMessageListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contact_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMessageListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageListResponse) ProtoMessage() {}

func (x *GetMessageListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contact_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageListResponse.ProtoReflect.Descriptor instead.
func (*GetMessageListResponse) Descriptor() ([]byte, []int) {
	return file_contact_proto_rawDescGZIP(), []int{6}
}

func (x *GetMessageListResponse) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

var File_contact_proto protoreflect.FileDescriptor

var file_contact_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x8d, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x54, 0x6f, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x6f, 0x49,
	0x64, 0x22, 0x4d, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x55, 0x73, 0x65, 0x72, 0x41, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x55,
	0x73, 0x65, 0x72, 0x41, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x42, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x55, 0x73, 0x65, 0x72, 0x42, 0x49, 0x64,
	0x22, 0x46, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5c, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x6f, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x6f, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x63, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x6f, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x6f, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x50,
	0x72, 0x65, 0x4d, 0x73, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x50, 0x72, 0x65, 0x4d, 0x73, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x46, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x32, 0xf5, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12,
	0x3e, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1d, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x57, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x47, 0x65,
	0x74, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e,
	0x47, 0x65, 0x74, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2e,
	0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contact_proto_rawDescOnce sync.Once
	file_contact_proto_rawDescData = file_contact_proto_rawDesc
)

func file_contact_proto_rawDescGZIP() []byte {
	file_contact_proto_rawDescOnce.Do(func() {
		file_contact_proto_rawDescData = protoimpl.X.CompressGZIP(file_contact_proto_rawDescData)
	})
	return file_contact_proto_rawDescData
}

var file_contact_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_contact_proto_goTypes = []interface{}{
	(*Empty)(nil),                    // 0: contact.Empty
	(*Message)(nil),                  // 1: contact.Message
	(*GetLatestMessageRequest)(nil),  // 2: contact.GetLatestMessageRequest
	(*GetLatestMessageResponse)(nil), // 3: contact.GetLatestMessageResponse
	(*CreateMessageRequest)(nil),     // 4: contact.CreateMessageRequest
	(*GetMessageListRequest)(nil),    // 5: contact.GetMessageListRequest
	(*GetMessageListResponse)(nil),   // 6: contact.GetMessageListResponse
}
var file_contact_proto_depIdxs = []int32{
	1, // 0: contact.GetLatestMessageResponse.Message:type_name -> contact.Message
	1, // 1: contact.GetMessageListResponse.Messages:type_name -> contact.Message
	4, // 2: contact.Contact.CreateMessage:input_type -> contact.CreateMessageRequest
	2, // 3: contact.Contact.GetLatestMessage:input_type -> contact.GetLatestMessageRequest
	5, // 4: contact.Contact.GetMessageList:input_type -> contact.GetMessageListRequest
	0, // 5: contact.Contact.CreateMessage:output_type -> contact.Empty
	3, // 6: contact.Contact.GetLatestMessage:output_type -> contact.GetLatestMessageResponse
	6, // 7: contact.Contact.GetMessageList:output_type -> contact.GetMessageListResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_contact_proto_init() }
func file_contact_proto_init() {
	if File_contact_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_contact_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_contact_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_contact_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLatestMessageRequest); i {
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
		file_contact_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLatestMessageResponse); i {
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
		file_contact_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMessageRequest); i {
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
		file_contact_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMessageListRequest); i {
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
		file_contact_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMessageListResponse); i {
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
			RawDescriptor: file_contact_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contact_proto_goTypes,
		DependencyIndexes: file_contact_proto_depIdxs,
		MessageInfos:      file_contact_proto_msgTypes,
	}.Build()
	File_contact_proto = out.File
	file_contact_proto_rawDesc = nil
	file_contact_proto_goTypes = nil
	file_contact_proto_depIdxs = nil
}
