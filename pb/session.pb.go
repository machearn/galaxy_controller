// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: session.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	UserId       int32                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	RefreshToken string                 `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	ClientIp     string                 `protobuf:"bytes,4,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
	UserAgent    string                 `protobuf:"bytes,5,opt,name=user_agent,json=userAgent,proto3" json:"user_agent,omitempty"`
	CreatedAt    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ExpiredAt    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=expired_at,json=expiredAt,proto3" json:"expired_at,omitempty"`
}

func (x *Session) Reset() {
	*x = Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_session_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_session_proto_rawDescGZIP(), []int{0}
}

func (x *Session) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Session) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Session) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *Session) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *Session) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

func (x *Session) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Session) GetExpiredAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiredAt
	}
	return nil
}

var File_session_proto protoreflect.FileDescriptor

var file_session_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x02, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1b,
	0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x41, 0x74,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x61, 0x63, 0x68, 0x65, 0x61, 0x72, 0x6e, 0x2f, 0x67, 0x61, 0x6c, 0x61, 0x78, 0x79, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_session_proto_rawDescOnce sync.Once
	file_session_proto_rawDescData = file_session_proto_rawDesc
)

func file_session_proto_rawDescGZIP() []byte {
	file_session_proto_rawDescOnce.Do(func() {
		file_session_proto_rawDescData = protoimpl.X.CompressGZIP(file_session_proto_rawDescData)
	})
	return file_session_proto_rawDescData
}

var file_session_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_session_proto_goTypes = []interface{}{
	(*Session)(nil),               // 0: pb.Session
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_session_proto_depIdxs = []int32{
	1, // 0: pb.Session.created_at:type_name -> google.protobuf.Timestamp
	1, // 1: pb.Session.expired_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_session_proto_init() }
func file_session_proto_init() {
	if File_session_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_session_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Session); i {
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
			RawDescriptor: file_session_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_session_proto_goTypes,
		DependencyIndexes: file_session_proto_depIdxs,
		MessageInfos:      file_session_proto_msgTypes,
	}.Build()
	File_session_proto = out.File
	file_session_proto_rawDesc = nil
	file_session_proto_goTypes = nil
	file_session_proto_depIdxs = nil
}
