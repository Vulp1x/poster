// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: parser/parser.proto

package parser

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

type SaveBotsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Bots  []*Bot `protobuf:"bytes,2,rep,name=bots,proto3" json:"bots,omitempty"`
}

func (x *SaveBotsRequest) Reset() {
	*x = SaveBotsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parser_parser_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBotsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBotsRequest) ProtoMessage() {}

func (x *SaveBotsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_parser_parser_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBotsRequest.ProtoReflect.Descriptor instead.
func (*SaveBotsRequest) Descriptor() ([]byte, []int) {
	return file_parser_parser_proto_rawDescGZIP(), []int{0}
}

func (x *SaveBotsRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *SaveBotsRequest) GetBots() []*Bot {
	if x != nil {
		return x.Bots
	}
	return nil
}

type Bot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// имя аккаунта в инстаграме
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// количество блогеров, которые проходят проверку по коду региона
	UserId int64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// количество блогеров, которые проходят проверку по коду региона
	SessionId string `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// прокси для бота
	Proxy *Proxy `protobuf:"bytes,4,opt,name=proxy,proto3" json:"proxy,omitempty"`
}

func (x *Bot) Reset() {
	*x = Bot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parser_parser_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bot) ProtoMessage() {}

func (x *Bot) ProtoReflect() protoreflect.Message {
	mi := &file_parser_parser_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bot.ProtoReflect.Descriptor instead.
func (*Bot) Descriptor() ([]byte, []int) {
	return file_parser_parser_proto_rawDescGZIP(), []int{1}
}

func (x *Bot) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Bot) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Bot) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *Bot) GetProxy() *Proxy {
	if x != nil {
		return x.Proxy
	}
	return nil
}

type Proxy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// имя аккаунта в инстаграме
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	// количество блогеров, которые проходят проверку по коду региона
	Port int32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// имя аккаунта в инстаграме
	Login string `protobuf:"bytes,3,opt,name=login,proto3" json:"login,omitempty"`
	// имя аккаунта в инстаграме
	Pass string `protobuf:"bytes,4,opt,name=pass,proto3" json:"pass,omitempty"`
}

func (x *Proxy) Reset() {
	*x = Proxy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parser_parser_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proxy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proxy) ProtoMessage() {}

func (x *Proxy) ProtoReflect() protoreflect.Message {
	mi := &file_parser_parser_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Proxy.ProtoReflect.Descriptor instead.
func (*Proxy) Descriptor() ([]byte, []int) {
	return file_parser_parser_proto_rawDescGZIP(), []int{2}
}

func (x *Proxy) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Proxy) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Proxy) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *Proxy) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

type SaveBotsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BotsSaved int32 `protobuf:"varint,1,opt,name=bots_saved,json=botsSaved,proto3" json:"bots_saved,omitempty"`
}

func (x *SaveBotsResponse) Reset() {
	*x = SaveBotsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_parser_parser_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBotsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBotsResponse) ProtoMessage() {}

func (x *SaveBotsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_parser_parser_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBotsResponse.ProtoReflect.Descriptor instead.
func (*SaveBotsResponse) Descriptor() ([]byte, []int) {
	return file_parser_parser_proto_rawDescGZIP(), []int{3}
}

func (x *SaveBotsResponse) GetBotsSaved() int32 {
	if x != nil {
		return x.BotsSaved
	}
	return 0
}

var File_parser_parser_proto protoreflect.FileDescriptor

var file_parser_parser_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x22, 0x48, 0x0a,
	0x0f, 0x53, 0x61, 0x76, 0x65, 0x42, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1f, 0x0a, 0x04, 0x62, 0x6f, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x42, 0x6f,
	0x74, 0x52, 0x04, 0x62, 0x6f, 0x74, 0x73, 0x22, 0x7e, 0x0a, 0x03, 0x42, 0x6f, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x52, 0x05, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x22, 0x59, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61,
	0x73, 0x73, 0x22, 0x31, 0x0a, 0x10, 0x53, 0x61, 0x76, 0x65, 0x42, 0x6f, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6f, 0x74, 0x73, 0x5f, 0x73,
	0x61, 0x76, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62, 0x6f, 0x74, 0x73,
	0x53, 0x61, 0x76, 0x65, 0x64, 0x32, 0x47, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x12,
	0x3d, 0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x42, 0x6f, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x70, 0x61,
	0x72, 0x73, 0x65, 0x72, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x42, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x53, 0x61,
	0x76, 0x65, 0x42, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x73,
	0x74, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x3b, 0x70, 0x61, 0x72, 0x73,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_parser_parser_proto_rawDescOnce sync.Once
	file_parser_parser_proto_rawDescData = file_parser_parser_proto_rawDesc
)

func file_parser_parser_proto_rawDescGZIP() []byte {
	file_parser_parser_proto_rawDescOnce.Do(func() {
		file_parser_parser_proto_rawDescData = protoimpl.X.CompressGZIP(file_parser_parser_proto_rawDescData)
	})
	return file_parser_parser_proto_rawDescData
}

var file_parser_parser_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_parser_parser_proto_goTypes = []interface{}{
	(*SaveBotsRequest)(nil),  // 0: parser.SaveBotsRequest
	(*Bot)(nil),              // 1: parser.Bot
	(*Proxy)(nil),            // 2: parser.Proxy
	(*SaveBotsResponse)(nil), // 3: parser.SaveBotsResponse
}
var file_parser_parser_proto_depIdxs = []int32{
	1, // 0: parser.SaveBotsRequest.bots:type_name -> parser.Bot
	2, // 1: parser.Bot.proxy:type_name -> parser.Proxy
	0, // 2: parser.Parser.SaveBots:input_type -> parser.SaveBotsRequest
	3, // 3: parser.Parser.SaveBots:output_type -> parser.SaveBotsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_parser_parser_proto_init() }
func file_parser_parser_proto_init() {
	if File_parser_parser_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_parser_parser_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBotsRequest); i {
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
		file_parser_parser_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bot); i {
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
		file_parser_parser_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Proxy); i {
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
		file_parser_parser_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBotsResponse); i {
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
			RawDescriptor: file_parser_parser_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_parser_parser_proto_goTypes,
		DependencyIndexes: file_parser_parser_proto_depIdxs,
		MessageInfos:      file_parser_parser_proto_msgTypes,
	}.Build()
	File_parser_parser_proto = out.File
	file_parser_parser_proto_rawDesc = nil
	file_parser_parser_proto_goTypes = nil
	file_parser_parser_proto_depIdxs = nil
}
