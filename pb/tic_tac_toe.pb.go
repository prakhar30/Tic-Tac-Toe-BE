// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: tic_tac_toe.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_tic_tac_toe_proto protoreflect.FileDescriptor

var file_tic_tac_toe_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x74, 0x69, 0x63, 0x5f, 0x74, 0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x74, 0x69, 0x63, 0x5f, 0x74, 0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65,
	0x1a, 0x15, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x72, 0x70, 0x63, 0x5f, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xaa, 0x01,
	0x0a, 0x09, 0x54, 0x69, 0x63, 0x54, 0x61, 0x63, 0x54, 0x6f, 0x65, 0x12, 0x4f, 0x0a, 0x0a, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x74, 0x69, 0x63, 0x5f,
	0x74, 0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x74, 0x69, 0x63, 0x5f,
	0x74, 0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x09,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x74, 0x69, 0x63, 0x5f,
	0x74, 0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x74, 0x69, 0x63, 0x5f, 0x74,
	0x61, 0x63, 0x5f, 0x74, 0x6f, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x6d, 0x61,
	0x69, 0x6e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var file_tic_tac_toe_proto_goTypes = []any{
	(*CreateUserRequest)(nil),  // 0: tic_tac_toe.CreateUserRequest
	(*LoginUserRequest)(nil),   // 1: tic_tac_toe.LoginUserRequest
	(*CreateUserResponse)(nil), // 2: tic_tac_toe.CreateUserResponse
	(*LoginUserResponse)(nil),  // 3: tic_tac_toe.LoginUserResponse
}
var file_tic_tac_toe_proto_depIdxs = []int32{
	0, // 0: tic_tac_toe.TicTacToe.CreateUser:input_type -> tic_tac_toe.CreateUserRequest
	1, // 1: tic_tac_toe.TicTacToe.LoginUser:input_type -> tic_tac_toe.LoginUserRequest
	2, // 2: tic_tac_toe.TicTacToe.CreateUser:output_type -> tic_tac_toe.CreateUserResponse
	3, // 3: tic_tac_toe.TicTacToe.LoginUser:output_type -> tic_tac_toe.LoginUserResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tic_tac_toe_proto_init() }
func file_tic_tac_toe_proto_init() {
	if File_tic_tac_toe_proto != nil {
		return
	}
	file_rpc_create_user_proto_init()
	file_rpc_login_user_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tic_tac_toe_proto_rawDesc), len(file_tic_tac_toe_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tic_tac_toe_proto_goTypes,
		DependencyIndexes: file_tic_tac_toe_proto_depIdxs,
	}.Build()
	File_tic_tac_toe_proto = out.File
	file_tic_tac_toe_proto_goTypes = nil
	file_tic_tac_toe_proto_depIdxs = nil
}
