// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: v1/error.proto

package v1

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

type ErrorCode int32

const (
	// Zero value is required and UNSPECIFIED is a convention.
	// Enum value names should be unique since they are global.
	//
	// ゼロは必須で UNSPECIFIED にする慣習。デフォルト値なので
	// enum 値名は unique である必要がある。プレフィクスをつけた方が良くはある
	ErrorCode_ERROR_CODE_UNSPECIFIED ErrorCode = 0
	// よく分からない時はこれを返せ
	ErrorCode_ERROR_CODE_INTERNAL ErrorCode = 1
	// ログイン
	ErrorCode_ERROR_CODE_USER_NOT_EXISTS ErrorCode = 20000
	// 認証
	ErrorCode_ERROR_CODE_TOKEN_INVALID ErrorCode = 20102
	ErrorCode_ERROR_CODE_TOKEN_EXPIRED ErrorCode = 20101
	// 認可
	ErrorCode_ERROR_CODE_PERMISSION_DENIED ErrorCode = 20201
	// その他
	ErrorCode_ERROR_CODE_INVALID_ARGUMENT ErrorCode = 20301
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0:     "ERROR_CODE_UNSPECIFIED",
		1:     "ERROR_CODE_INTERNAL",
		20000: "ERROR_CODE_USER_NOT_EXISTS",
		20102: "ERROR_CODE_TOKEN_INVALID",
		20101: "ERROR_CODE_TOKEN_EXPIRED",
		20201: "ERROR_CODE_PERMISSION_DENIED",
		20301: "ERROR_CODE_INVALID_ARGUMENT",
	}
	ErrorCode_value = map[string]int32{
		"ERROR_CODE_UNSPECIFIED":       0,
		"ERROR_CODE_INTERNAL":          1,
		"ERROR_CODE_USER_NOT_EXISTS":   20000,
		"ERROR_CODE_TOKEN_INVALID":     20102,
		"ERROR_CODE_TOKEN_EXPIRED":     20101,
		"ERROR_CODE_PERMISSION_DENIED": 20201,
		"ERROR_CODE_INVALID_ARGUMENT":  20301,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_error_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_v1_error_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_v1_error_proto_rawDescGZIP(), []int{0}
}

var File_v1_error_proto protoreflect.FileDescriptor

var file_v1_error_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x76, 0x31, 0x2a, 0xe9, 0x01, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17,
	0x0a, 0x13, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x20, 0x0a, 0x1a, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x45,
	0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0xa0, 0x9c, 0x01, 0x12, 0x1e, 0x0a, 0x18, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x5f, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x86, 0x9d, 0x01, 0x12, 0x1e, 0x0a, 0x18, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x5f, 0x45, 0x58,
	0x50, 0x49, 0x52, 0x45, 0x44, 0x10, 0x85, 0x9d, 0x01, 0x12, 0x22, 0x0a, 0x1c, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49,
	0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44, 0x10, 0xe9, 0x9d, 0x01, 0x12, 0x21, 0x0a,
	0x1b, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41,
	0x4c, 0x49, 0x44, 0x5f, 0x41, 0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0xcd, 0x9e, 0x01,
	0x42, 0x4a, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0c, 0x61, 0x70, 0x70, 0x2f, 0x70, 0x62,
	0x67, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x56, 0x58, 0x58, 0xaa, 0x02, 0x02, 0x56,
	0x31, 0xca, 0x02, 0x02, 0x56, 0x31, 0xe2, 0x02, 0x0e, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x02, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_error_proto_rawDescOnce sync.Once
	file_v1_error_proto_rawDescData = file_v1_error_proto_rawDesc
)

func file_v1_error_proto_rawDescGZIP() []byte {
	file_v1_error_proto_rawDescOnce.Do(func() {
		file_v1_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_error_proto_rawDescData)
	})
	return file_v1_error_proto_rawDescData
}

var file_v1_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v1_error_proto_goTypes = []interface{}{
	(ErrorCode)(0), // 0: v1.ErrorCode
}
var file_v1_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_error_proto_init() }
func file_v1_error_proto_init() {
	if File_v1_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_error_proto_goTypes,
		DependencyIndexes: file_v1_error_proto_depIdxs,
		EnumInfos:         file_v1_error_proto_enumTypes,
	}.Build()
	File_v1_error_proto = out.File
	file_v1_error_proto_rawDesc = nil
	file_v1_error_proto_goTypes = nil
	file_v1_error_proto_depIdxs = nil
}
