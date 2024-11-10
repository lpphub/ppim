// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: message.proto

package message_pb

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
	MsgType_SEND        MsgType = 2
	MsgType_SEND_ACK    MsgType = 3
	MsgType_RECEIVE     MsgType = 4
	MsgType_RECEIVE_ACK MsgType = 5
	MsgType_PING        MsgType = 6
	MsgType_PONG        MsgType = 7
	MsgType_DISCONNECT  MsgType = 8
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "UNKNOWN",
		1: "CONNECT",
		2: "SEND",
		3: "SEND_ACK",
		4: "RECEIVE",
		5: "RECEIVE_ACK",
		6: "PING",
		7: "PONG",
		8: "DISCONNECT",
	}
	MsgType_value = map[string]int32{
		"UNKNOWN":     0,
		"CONNECT":     1,
		"SEND":        2,
		"SEND_ACK":    3,
		"RECEIVE":     4,
		"RECEIVE_ACK": 5,
		"PING":        6,
		"PONG":        7,
		"DISCONNECT":  8,
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
	return file_message_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType MsgType `protobuf:"varint,1,opt,name=msgType,proto3,enum=message.MsgType" json:"msgType,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*Message_ConnectPacket
	//	*Message_SendPacket
	Payload isMessage_Payload `protobuf_oneof:"payload"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
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
	return file_message_proto_rawDescGZIP(), []int{0}
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

func (x *Message) GetSendPacket() *SendPacket {
	if x, ok := x.GetPayload().(*Message_SendPacket); ok {
		return x.SendPacket
	}
	return nil
}

type isMessage_Payload interface {
	isMessage_Payload()
}

type Message_ConnectPacket struct {
	ConnectPacket *ConnectPacket `protobuf:"bytes,2,opt,name=connectPacket,proto3,oneof"`
}

type Message_SendPacket struct {
	SendPacket *SendPacket `protobuf:"bytes,3,opt,name=sendPacket,proto3,oneof"`
}

func (*Message_ConnectPacket) isMessage_Payload() {}

func (*Message_SendPacket) isMessage_Payload() {}

type ConnectPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ConnectPacket) Reset() {
	*x = ConnectPacket{}
	mi := &file_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConnectPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectPacket) ProtoMessage() {}

func (x *ConnectPacket) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
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
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectPacket) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ConnectPacket) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SendPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *SendPacket) Reset() {
	*x = SendPacket{}
	mi := &file_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendPacket) ProtoMessage() {}

func (x *SendPacket) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
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
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *SendPacket) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xb7, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x3e, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48,
	0x00, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x35, 0x0a, 0x0a, 0x73, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x73, 0x65, 0x6e,
	0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x22, 0x3d, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x26, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2a, 0x7d, 0x0a, 0x07, 0x4d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x01, 0x12, 0x08,
	0x0a, 0x04, 0x53, 0x45, 0x4e, 0x44, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x45, 0x4e, 0x44,
	0x5f, 0x41, 0x43, 0x4b, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x43, 0x45, 0x49, 0x56,
	0x45, 0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x43, 0x45, 0x49, 0x56, 0x45, 0x5f, 0x41,
	0x43, 0x4b, 0x10, 0x05, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x4f, 0x4e, 0x47, 0x10, 0x07, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x49, 0x53, 0x43,
	0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x08, 0x42, 0x19, 0x5a, 0x17, 0x2e, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x62, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_message_proto_goTypes = []any{
	(MsgType)(0),          // 0: message.MsgType
	(*Message)(nil),       // 1: message.Message
	(*ConnectPacket)(nil), // 2: message.ConnectPacket
	(*SendPacket)(nil),    // 3: message.SendPacket
}
var file_message_proto_depIdxs = []int32{
	0, // 0: message.Message.msgType:type_name -> message.MsgType
	2, // 1: message.Message.connectPacket:type_name -> message.ConnectPacket
	3, // 2: message.Message.sendPacket:type_name -> message.SendPacket
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	file_message_proto_msgTypes[0].OneofWrappers = []any{
		(*Message_ConnectPacket)(nil),
		(*Message_SendPacket)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		EnumInfos:         file_message_proto_enumTypes,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
