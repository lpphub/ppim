// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: protocol.proto

package protocol

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

type MsgType int32

const (
	MsgType_UNKNOWN     MsgType = 0
	MsgType_CONNECT     MsgType = 1
	MsgType_CONNECT_ACK MsgType = 2
	MsgType_SEND        MsgType = 3 // 发送消息 (c2s)
	MsgType_SEND_ACK    MsgType = 4 // 发送消息 ack (s2c)
	MsgType_RECEIVE     MsgType = 5 // 接收消息 (s2c)
	MsgType_RECEIVE_ACK MsgType = 6 // 接收消息 ack (c2s)
	MsgType_PING        MsgType = 7
	MsgType_PONG        MsgType = 8
	MsgType_DISCONNECT  MsgType = 9
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "UNKNOWN",
		1: "CONNECT",
		2: "CONNECT_ACK",
		3: "SEND",
		4: "SEND_ACK",
		5: "RECEIVE",
		6: "RECEIVE_ACK",
		7: "PING",
		8: "PONG",
		9: "DISCONNECT",
	}
	MsgType_value = map[string]int32{
		"UNKNOWN":     0,
		"CONNECT":     1,
		"CONNECT_ACK": 2,
		"SEND":        3,
		"SEND_ACK":    4,
		"RECEIVE":     5,
		"RECEIVE_ACK": 6,
		"PING":        7,
		"PONG":        8,
		"DISCONNECT":  9,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_protocol_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_protocol_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType MsgType `protobuf:"varint,1,opt,name=msgType,proto3,enum=message.MsgType" json:"msgType,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*Message_ConnectPacket
	//	*Message_ConnectAckPacket
	//	*Message_PingPacket
	//	*Message_PongPacket
	//	*Message_SendPacket
	//	*Message_SendAckPacket
	//	*Message_ReceivePacket
	//	*Message_ReceiveAckPacket
	Payload isMessage_Payload `protobuf_oneof:"payload"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_protocol_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[0]
	if x != nil {
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
	return file_protocol_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetMsgType() MsgType {
	if x != nil {
		return x.MsgType
	}
	return MsgType_UNKNOWN
}

func (m *Message) GetPayload() isMessage_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *Message) GetConnectPacket() *ConnectPacket {
	if x, ok := x.GetPayload().(*Message_ConnectPacket); ok {
		return x.ConnectPacket
	}
	return nil
}

func (x *Message) GetConnectAckPacket() *ConnectAckPacket {
	if x, ok := x.GetPayload().(*Message_ConnectAckPacket); ok {
		return x.ConnectAckPacket
	}
	return nil
}

func (x *Message) GetPingPacket() *PingPacket {
	if x, ok := x.GetPayload().(*Message_PingPacket); ok {
		return x.PingPacket
	}
	return nil
}

func (x *Message) GetPongPacket() *PongPacket {
	if x, ok := x.GetPayload().(*Message_PongPacket); ok {
		return x.PongPacket
	}
	return nil
}

func (x *Message) GetSendPacket() *SendPacket {
	if x, ok := x.GetPayload().(*Message_SendPacket); ok {
		return x.SendPacket
	}
	return nil
}

func (x *Message) GetSendAckPacket() *SendAckPacket {
	if x, ok := x.GetPayload().(*Message_SendAckPacket); ok {
		return x.SendAckPacket
	}
	return nil
}

func (x *Message) GetReceivePacket() *ReceivePacket {
	if x, ok := x.GetPayload().(*Message_ReceivePacket); ok {
		return x.ReceivePacket
	}
	return nil
}

func (x *Message) GetReceiveAckPacket() *ReceiveAckPacket {
	if x, ok := x.GetPayload().(*Message_ReceiveAckPacket); ok {
		return x.ReceiveAckPacket
	}
	return nil
}

type isMessage_Payload interface {
	isMessage_Payload()
}

type Message_ConnectPacket struct {
	ConnectPacket *ConnectPacket `protobuf:"bytes,2,opt,name=connectPacket,proto3,oneof"`
}

type Message_ConnectAckPacket struct {
	ConnectAckPacket *ConnectAckPacket `protobuf:"bytes,3,opt,name=connectAckPacket,proto3,oneof"`
}

type Message_PingPacket struct {
	PingPacket *PingPacket `protobuf:"bytes,4,opt,name=pingPacket,proto3,oneof"`
}

type Message_PongPacket struct {
	PongPacket *PongPacket `protobuf:"bytes,5,opt,name=pongPacket,proto3,oneof"`
}

type Message_SendPacket struct {
	SendPacket *SendPacket `protobuf:"bytes,6,opt,name=sendPacket,proto3,oneof"`
}

type Message_SendAckPacket struct {
	SendAckPacket *SendAckPacket `protobuf:"bytes,7,opt,name=sendAckPacket,proto3,oneof"`
}

type Message_ReceivePacket struct {
	ReceivePacket *ReceivePacket `protobuf:"bytes,8,opt,name=receivePacket,proto3,oneof"`
}

type Message_ReceiveAckPacket struct {
	ReceiveAckPacket *ReceiveAckPacket `protobuf:"bytes,9,opt,name=receiveAckPacket,proto3,oneof"`
}

func (*Message_ConnectPacket) isMessage_Payload() {}

func (*Message_ConnectAckPacket) isMessage_Payload() {}

func (*Message_PingPacket) isMessage_Payload() {}

func (*Message_PongPacket) isMessage_Payload() {}

func (*Message_SendPacket) isMessage_Payload() {}

func (*Message_SendAckPacket) isMessage_Payload() {}

func (*Message_ReceivePacket) isMessage_Payload() {}

func (*Message_ReceiveAckPacket) isMessage_Payload() {}

type ConnectPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid   string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Did   string `protobuf:"bytes,3,opt,name=did,proto3" json:"did,omitempty"`
}

func (x *ConnectPacket) Reset() {
	*x = ConnectPacket{}
	mi := &file_protocol_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConnectPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectPacket) ProtoMessage() {}

func (x *ConnectPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectPacket.ProtoReflect.Descriptor instead.
func (*ConnectPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectPacket) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *ConnectPacket) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ConnectPacket) GetDid() string {
	if x != nil {
		return x.Did
	}
	return ""
}

type ConnectAckPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *ConnectAckPacket) Reset() {
	*x = ConnectAckPacket{}
	mi := &file_protocol_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConnectAckPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectAckPacket) ProtoMessage() {}

func (x *ConnectAckPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectAckPacket.ProtoReflect.Descriptor instead.
func (*ConnectAckPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{2}
}

func (x *ConnectAckPacket) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type PingPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingPacket) Reset() {
	*x = PingPacket{}
	mi := &file_protocol_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingPacket) ProtoMessage() {}

func (x *PingPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingPacket.ProtoReflect.Descriptor instead.
func (*PingPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{3}
}

type PongPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PongPacket) Reset() {
	*x = PongPacket{}
	mi := &file_protocol_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PongPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PongPacket) ProtoMessage() {}

func (x *PongPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PongPacket.ProtoReflect.Descriptor instead.
func (*PongPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{4}
}

type SendPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationType string   `protobuf:"bytes,1,opt,name=conversationType,proto3" json:"conversationType,omitempty"`
	ToID             string   `protobuf:"bytes,2,opt,name=toID,proto3" json:"toID,omitempty"`
	Payload          *Payload `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *SendPacket) Reset() {
	*x = SendPacket{}
	mi := &file_protocol_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendPacket) ProtoMessage() {}

func (x *SendPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendPacket.ProtoReflect.Descriptor instead.
func (*SendPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{5}
}

func (x *SendPacket) GetConversationType() string {
	if x != nil {
		return x.ConversationType
	}
	return ""
}

func (x *SendPacket) GetToID() string {
	if x != nil {
		return x.ToID
	}
	return ""
}

func (x *SendPacket) GetPayload() *Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType int32  `protobuf:"varint,1,opt,name=msgType,proto3" json:"msgType,omitempty"`
	MsgNo   string `protobuf:"bytes,2,opt,name=msgNo,proto3" json:"msgNo,omitempty"`
	Content string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Payload) Reset() {
	*x = Payload{}
	mi := &file_protocol_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payload) ProtoMessage() {}

func (x *Payload) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payload.ProtoReflect.Descriptor instead.
func (*Payload) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{6}
}

func (x *Payload) GetMsgType() int32 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *Payload) GetMsgNo() string {
	if x != nil {
		return x.MsgNo
	}
	return ""
}

func (x *Payload) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type SendAckPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	MsgNo  string `protobuf:"bytes,2,opt,name=msgNo,proto3" json:"msgNo,omitempty"`
	MsgId  string `protobuf:"bytes,3,opt,name=msgId,proto3" json:"msgId,omitempty"`
	MsgSeq uint64 `protobuf:"varint,4,opt,name=msgSeq,proto3" json:"msgSeq,omitempty"`
}

func (x *SendAckPacket) Reset() {
	*x = SendAckPacket{}
	mi := &file_protocol_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendAckPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendAckPacket) ProtoMessage() {}

func (x *SendAckPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendAckPacket.ProtoReflect.Descriptor instead.
func (*SendAckPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{7}
}

func (x *SendAckPacket) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *SendAckPacket) GetMsgNo() string {
	if x != nil {
		return x.MsgNo
	}
	return ""
}

func (x *SendAckPacket) GetMsgId() string {
	if x != nil {
		return x.MsgId
	}
	return ""
}

func (x *SendAckPacket) GetMsgSeq() uint64 {
	if x != nil {
		return x.MsgSeq
	}
	return 0
}

type ReceivePacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ReceivePacket) Reset() {
	*x = ReceivePacket{}
	mi := &file_protocol_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReceivePacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceivePacket) ProtoMessage() {}

func (x *ReceivePacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceivePacket.ProtoReflect.Descriptor instead.
func (*ReceivePacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{8}
}

func (x *ReceivePacket) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type ReceiveAckPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	MsgId  string `protobuf:"bytes,2,opt,name=msgId,proto3" json:"msgId,omitempty"`
	MsgSeq uint64 `protobuf:"varint,3,opt,name=msgSeq,proto3" json:"msgSeq,omitempty"`
}

func (x *ReceiveAckPacket) Reset() {
	*x = ReceiveAckPacket{}
	mi := &file_protocol_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReceiveAckPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveAckPacket) ProtoMessage() {}

func (x *ReceiveAckPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveAckPacket.ProtoReflect.Descriptor instead.
func (*ReceiveAckPacket) Descriptor() ([]byte, []int) {
	return file_protocol_proto_rawDescGZIP(), []int{9}
}

func (x *ReceiveAckPacket) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ReceiveAckPacket) GetMsgId() string {
	if x != nil {
		return x.MsgId
	}
	return ""
}

func (x *ReceiveAckPacket) GetMsgSeq() uint64 {
	if x != nil {
		return x.MsgSeq
	}
	return 0
}

var File_protocol_proto protoreflect.FileDescriptor

var file_protocol_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xb7, 0x04, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x48, 0x00, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x12, 0x47, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x41, 0x63, 0x6b, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x41, 0x63, 0x6b,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x70, 0x69,
	0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x12, 0x35, 0x0a, 0x0a, 0x70, 0x6f, 0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x50, 0x6f, 0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x6f,
	0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x73, 0x65, 0x6e, 0x64,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x48, 0x00, 0x52, 0x0a, 0x73, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x3e, 0x0a, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00,
	0x52, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x3e, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00,
	0x52, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x47, 0x0a, 0x10, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x41, 0x63, 0x6b, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x10, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x41,
	0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x49, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x64, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x69, 0x64, 0x22, 0x26,
	0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x0c, 0x0a, 0x0a, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x22, 0x0c, 0x0a, 0x0a, 0x50, 0x6f, 0x6e, 0x67, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x22, 0x78, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x2a, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x6f, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x6f, 0x49, 0x44,
	0x12, 0x2a, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x53, 0x0a, 0x07,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x73, 0x67, 0x4e, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6d, 0x73, 0x67, 0x4e, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x22, 0x67, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x73, 0x67, 0x4e, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x4e, 0x6f, 0x12, 0x14, 0x0a, 0x05,
	0x6d, 0x73, 0x67, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x73, 0x67,
	0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x53, 0x65, 0x71, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x6d, 0x73, 0x67, 0x53, 0x65, 0x71, 0x22, 0x29, 0x0a, 0x0d, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x54, 0x0a, 0x10, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x41, 0x63, 0x6b, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x73,
	0x67, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x53, 0x65, 0x71, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6d, 0x73, 0x67, 0x53, 0x65, 0x71, 0x2a, 0x8e, 0x01, 0x0a, 0x07,
	0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f,
	0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10,
	0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x5f, 0x41, 0x43, 0x4b,
	0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x45, 0x4e, 0x44, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08,
	0x53, 0x45, 0x4e, 0x44, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45,
	0x43, 0x45, 0x49, 0x56, 0x45, 0x10, 0x05, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x43, 0x45, 0x49,
	0x56, 0x45, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47,
	0x10, 0x07, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4f, 0x4e, 0x47, 0x10, 0x08, 0x12, 0x0e, 0x0a, 0x0a,
	0x44, 0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x09, 0x42, 0x15, 0x5a, 0x13,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protocol_proto_rawDescOnce sync.Once
	file_protocol_proto_rawDescData = file_protocol_proto_rawDesc
)

func file_protocol_proto_rawDescGZIP() []byte {
	file_protocol_proto_rawDescOnce.Do(func() {
		file_protocol_proto_rawDescData = protoimpl.X.CompressGZIP(file_protocol_proto_rawDescData)
	})
	return file_protocol_proto_rawDescData
}

var file_protocol_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protocol_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_protocol_proto_goTypes = []any{
	(MsgType)(0),             // 0: message.MsgType
	(*Message)(nil),          // 1: message.Message
	(*ConnectPacket)(nil),    // 2: message.ConnectPacket
	(*ConnectAckPacket)(nil), // 3: message.ConnectAckPacket
	(*PingPacket)(nil),       // 4: message.PingPacket
	(*PongPacket)(nil),       // 5: message.PongPacket
	(*SendPacket)(nil),       // 6: message.SendPacket
	(*Payload)(nil),          // 7: message.Payload
	(*SendAckPacket)(nil),    // 8: message.SendAckPacket
	(*ReceivePacket)(nil),    // 9: message.ReceivePacket
	(*ReceiveAckPacket)(nil), // 10: message.ReceiveAckPacket
}
var file_protocol_proto_depIdxs = []int32{
	0,  // 0: message.Message.msgType:type_name -> message.MsgType
	2,  // 1: message.Message.connectPacket:type_name -> message.ConnectPacket
	3,  // 2: message.Message.connectAckPacket:type_name -> message.ConnectAckPacket
	4,  // 3: message.Message.pingPacket:type_name -> message.PingPacket
	5,  // 4: message.Message.pongPacket:type_name -> message.PongPacket
	6,  // 5: message.Message.sendPacket:type_name -> message.SendPacket
	8,  // 6: message.Message.sendAckPacket:type_name -> message.SendAckPacket
	9,  // 7: message.Message.receivePacket:type_name -> message.ReceivePacket
	10, // 8: message.Message.receiveAckPacket:type_name -> message.ReceiveAckPacket
	7,  // 9: message.SendPacket.payload:type_name -> message.Payload
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_protocol_proto_init() }
func file_protocol_proto_init() {
	if File_protocol_proto != nil {
		return
	}
	file_protocol_proto_msgTypes[0].OneofWrappers = []any{
		(*Message_ConnectPacket)(nil),
		(*Message_ConnectAckPacket)(nil),
		(*Message_PingPacket)(nil),
		(*Message_PongPacket)(nil),
		(*Message_SendPacket)(nil),
		(*Message_SendAckPacket)(nil),
		(*Message_ReceivePacket)(nil),
		(*Message_ReceiveAckPacket)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protocol_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protocol_proto_goTypes,
		DependencyIndexes: file_protocol_proto_depIdxs,
		EnumInfos:         file_protocol_proto_enumTypes,
		MessageInfos:      file_protocol_proto_msgTypes,
	}.Build()
	File_protocol_proto = out.File
	file_protocol_proto_rawDesc = nil
	file_protocol_proto_goTypes = nil
	file_protocol_proto_depIdxs = nil
}
