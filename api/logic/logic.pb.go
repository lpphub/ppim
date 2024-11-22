// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: logic.proto

package logic

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

type AuthReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid   string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Did   string `protobuf:"bytes,2,opt,name=did,proto3" json:"did,omitempty"`
	Token string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AuthReq) Reset() {
	*x = AuthReq{}
	mi := &file_logic_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthReq) ProtoMessage() {}

func (x *AuthReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthReq.ProtoReflect.Descriptor instead.
func (*AuthReq) Descriptor() ([]byte, []int) {
	return file_logic_proto_rawDescGZIP(), []int{0}
}

func (x *AuthReq) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *AuthReq) GetDid() string {
	if x != nil {
		return x.Did
	}
	return ""
}

func (x *AuthReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AuthResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool    `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg *string `protobuf:"bytes,2,opt,name=msg,proto3,oneof" json:"msg,omitempty"`
}

func (x *AuthResp) Reset() {
	*x = AuthResp{}
	mi := &file_logic_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResp) ProtoMessage() {}

func (x *AuthResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResp.ProtoReflect.Descriptor instead.
func (*AuthResp) Descriptor() ([]byte, []int) {
	return file_logic_proto_rawDescGZIP(), []int{1}
}

func (x *AuthResp) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *AuthResp) GetMsg() string {
	if x != nil && x.Msg != nil {
		return *x.Msg
	}
	return ""
}

type OnlineReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid   string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Did   string `protobuf:"bytes,2,opt,name=did,proto3" json:"did,omitempty"`
	Ip    string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	Topic string `protobuf:"bytes,4,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *OnlineReq) Reset() {
	*x = OnlineReq{}
	mi := &file_logic_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OnlineReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineReq) ProtoMessage() {}

func (x *OnlineReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineReq.ProtoReflect.Descriptor instead.
func (*OnlineReq) Descriptor() ([]byte, []int) {
	return file_logic_proto_rawDescGZIP(), []int{2}
}

func (x *OnlineReq) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *OnlineReq) GetDid() string {
	if x != nil {
		return x.Did
	}
	return ""
}

func (x *OnlineReq) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *OnlineReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type OnlineResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool    `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg *string `protobuf:"bytes,2,opt,name=msg,proto3,oneof" json:"msg,omitempty"`
}

func (x *OnlineResp) Reset() {
	*x = OnlineResp{}
	mi := &file_logic_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OnlineResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineResp) ProtoMessage() {}

func (x *OnlineResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineResp.ProtoReflect.Descriptor instead.
func (*OnlineResp) Descriptor() ([]byte, []int) {
	return file_logic_proto_rawDescGZIP(), []int{3}
}

func (x *OnlineResp) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *OnlineResp) GetMsg() string {
	if x != nil && x.Msg != nil {
		return *x.Msg
	}
	return ""
}

var File_logic_proto protoreflect.FileDescriptor

var file_logic_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6c,
	0x6f, 0x67, 0x69, 0x63, 0x22, 0x43, 0x0a, 0x07, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x64, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x39, 0x0a, 0x08, 0x41, 0x75, 0x74,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x15, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x6d, 0x73, 0x67, 0x22, 0x55, 0x0a, 0x09, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x64, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x3b, 0x0a, 0x0a, 0x4f,
	0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x15, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x88, 0x01, 0x01,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x73, 0x67, 0x32, 0x2f, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x27, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x32, 0x6c, 0x0a, 0x06, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x10, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65,
	0x71, 0x1a, 0x11, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x31, 0x0a, 0x0a, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x10, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x6c, 0x6f, 0x67,
	0x69, 0x63, 0x3b, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_logic_proto_rawDescOnce sync.Once
	file_logic_proto_rawDescData = file_logic_proto_rawDesc
)

func file_logic_proto_rawDescGZIP() []byte {
	file_logic_proto_rawDescOnce.Do(func() {
		file_logic_proto_rawDescData = protoimpl.X.CompressGZIP(file_logic_proto_rawDescData)
	})
	return file_logic_proto_rawDescData
}

var file_logic_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_logic_proto_goTypes = []any{
	(*AuthReq)(nil),    // 0: logic.AuthReq
	(*AuthResp)(nil),   // 1: logic.AuthResp
	(*OnlineReq)(nil),  // 2: logic.OnlineReq
	(*OnlineResp)(nil), // 3: logic.OnlineResp
}
var file_logic_proto_depIdxs = []int32{
	0, // 0: logic.Auth.Auth:input_type -> logic.AuthReq
	2, // 1: logic.Online.Register:input_type -> logic.OnlineReq
	2, // 2: logic.Online.Unregister:input_type -> logic.OnlineReq
	1, // 3: logic.Auth.Auth:output_type -> logic.AuthResp
	3, // 4: logic.Online.Register:output_type -> logic.OnlineResp
	3, // 5: logic.Online.Unregister:output_type -> logic.OnlineResp
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_logic_proto_init() }
func file_logic_proto_init() {
	if File_logic_proto != nil {
		return
	}
	file_logic_proto_msgTypes[1].OneofWrappers = []any{}
	file_logic_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_logic_proto_goTypes,
		DependencyIndexes: file_logic_proto_depIdxs,
		MessageInfos:      file_logic_proto_msgTypes,
	}.Build()
	File_logic_proto = out.File
	file_logic_proto_rawDesc = nil
	file_logic_proto_goTypes = nil
	file_logic_proto_depIdxs = nil
}
