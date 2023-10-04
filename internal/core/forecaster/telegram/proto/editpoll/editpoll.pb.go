// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1-devel
// 	protoc        v3.6.1
// source: internal/core/forecaster/telegram/proto/editpoll/editpoll.proto

package editpoll

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

type FieldToEdit int32

const (
	FieldToEdit_UNDEFINED   FieldToEdit = 0
	FieldToEdit_TITLE       FieldToEdit = 1
	FieldToEdit_DESCRIPTION FieldToEdit = 2
	FieldToEdit_START_DATE  FieldToEdit = 3
	FieldToEdit_FINISH_DATE FieldToEdit = 4
)

// Enum value maps for FieldToEdit.
var (
	FieldToEdit_name = map[int32]string{
		0: "UNDEFINED",
		1: "TITLE",
		2: "DESCRIPTION",
		3: "START_DATE",
		4: "FINISH_DATE",
	}
	FieldToEdit_value = map[string]int32{
		"UNDEFINED":   0,
		"TITLE":       1,
		"DESCRIPTION": 2,
		"START_DATE":  3,
		"FINISH_DATE": 4,
	}
)

func (x FieldToEdit) Enum() *FieldToEdit {
	p := new(FieldToEdit)
	*p = x
	return p
}

func (x FieldToEdit) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldToEdit) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_enumTypes[0].Descriptor()
}

func (FieldToEdit) Type() protoreflect.EnumType {
	return &file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_enumTypes[0]
}

func (x FieldToEdit) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *FieldToEdit) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = FieldToEdit(num)
	return nil
}

// Deprecated: Use FieldToEdit.Descriptor instead.
func (FieldToEdit) EnumDescriptor() ([]byte, []int) {
	return file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescGZIP(), []int{0}
}

type EditPoll struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PollId      *int32       `protobuf:"varint,1,req,name=poll_id,json=pollId" json:"poll_id,omitempty"`
	FieldToEdit *FieldToEdit `protobuf:"varint,2,req,name=field_to_edit,json=fieldToEdit,enum=editpoll.FieldToEdit" json:"field_to_edit,omitempty"`
}

func (x *EditPoll) Reset() {
	*x = EditPoll{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditPoll) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditPoll) ProtoMessage() {}

func (x *EditPoll) ProtoReflect() protoreflect.Message {
	mi := &file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditPoll.ProtoReflect.Descriptor instead.
func (*EditPoll) Descriptor() ([]byte, []int) {
	return file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescGZIP(), []int{0}
}

func (x *EditPoll) GetPollId() int32 {
	if x != nil && x.PollId != nil {
		return *x.PollId
	}
	return 0
}

func (x *EditPoll) GetFieldToEdit() FieldToEdit {
	if x != nil && x.FieldToEdit != nil {
		return *x.FieldToEdit
	}
	return FieldToEdit_UNDEFINED
}

var File_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto protoreflect.FileDescriptor

var file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDesc = []byte{
	0x0a, 0x3f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x67,
	0x72, 0x61, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x70, 0x6f,
	0x6c, 0x6c, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x70, 0x6f, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x65, 0x64, 0x69, 0x74, 0x70, 0x6f, 0x6c, 0x6c, 0x22, 0x5e, 0x0a, 0x08, 0x45,
	0x64, 0x69, 0x74, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x6c, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x49, 0x64,
	0x12, 0x39, 0x0a, 0x0d, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x74, 0x6f, 0x5f, 0x65, 0x64, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x70, 0x6f,
	0x6c, 0x6c, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x6f, 0x45, 0x64, 0x69, 0x74, 0x52, 0x0b,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x6f, 0x45, 0x64, 0x69, 0x74, 0x2a, 0x59, 0x0a, 0x0b, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x54, 0x6f, 0x45, 0x64, 0x69, 0x74, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4e,
	0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x49, 0x54,
	0x4c, 0x45, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x45, 0x53, 0x43, 0x52, 0x49, 0x50, 0x54,
	0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x52, 0x54, 0x5f, 0x44,
	0x41, 0x54, 0x45, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x5f,
	0x44, 0x41, 0x54, 0x45, 0x10, 0x04, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x65, 0x64, 0x69, 0x74,
	0x70, 0x6f, 0x6c, 0x6c,
}

var (
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescOnce sync.Once
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescData = file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDesc
)

func file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescGZIP() []byte {
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescOnce.Do(func() {
		file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescData)
	})
	return file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDescData
}

var file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_goTypes = []interface{}{
	(FieldToEdit)(0), // 0: editpoll.FieldToEdit
	(*EditPoll)(nil), // 1: editpoll.EditPoll
}
var file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_depIdxs = []int32{
	0, // 0: editpoll.EditPoll.field_to_edit:type_name -> editpoll.FieldToEdit
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_init() }
func file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_init() {
	if File_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditPoll); i {
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
			RawDescriptor: file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_goTypes,
		DependencyIndexes: file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_depIdxs,
		EnumInfos:         file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_enumTypes,
		MessageInfos:      file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_msgTypes,
	}.Build()
	File_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto = out.File
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_rawDesc = nil
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_goTypes = nil
	file_internal_core_forecaster_telegram_proto_editpoll_editpoll_proto_depIdxs = nil
}
