// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1-devel
// 	protoc        v3.6.1
// source: internal/core/forecaster/telegram/proto/editfield/editfield.proto

package editfield

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

type Field int32

const (
	Field_UNDEFINED   Field = 0
	Field_TITLE       Field = 1
	Field_DESCRIPTION Field = 2
	Field_START_DATE  Field = 3
	Field_FINISH_DATE Field = 4
)

// Enum value maps for Field.
var (
	Field_name = map[int32]string{
		0: "UNDEFINED",
		1: "TITLE",
		2: "DESCRIPTION",
		3: "START_DATE",
		4: "FINISH_DATE",
	}
	Field_value = map[string]int32{
		"UNDEFINED":   0,
		"TITLE":       1,
		"DESCRIPTION": 2,
		"START_DATE":  3,
		"FINISH_DATE": 4,
	}
)

func (x Field) Enum() *Field {
	p := new(Field)
	*p = x
	return p
}

func (x Field) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Field) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_enumTypes[0].Descriptor()
}

func (Field) Type() protoreflect.EnumType {
	return &file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_enumTypes[0]
}

func (x Field) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Field) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Field(num)
	return nil
}

// Deprecated: Use Field.Descriptor instead.
func (Field) EnumDescriptor() ([]byte, []int) {
	return file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescGZIP(), []int{0}
}

type EditField struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PollId              *int32 `protobuf:"varint,1,req,name=poll_id,json=pollId" json:"poll_id,omitempty"`
	Field               *Field `protobuf:"varint,3,req,name=field,enum=editfield.Field" json:"field,omitempty"`
	ReferrerMyPollsPage *int32 `protobuf:"varint,4,opt,name=referrer_my_polls_page,json=referrerMyPollsPage" json:"referrer_my_polls_page,omitempty"`
}

func (x *EditField) Reset() {
	*x = EditField{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditField) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditField) ProtoMessage() {}

func (x *EditField) ProtoReflect() protoreflect.Message {
	mi := &file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditField.ProtoReflect.Descriptor instead.
func (*EditField) Descriptor() ([]byte, []int) {
	return file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescGZIP(), []int{0}
}

func (x *EditField) GetPollId() int32 {
	if x != nil && x.PollId != nil {
		return *x.PollId
	}
	return 0
}

func (x *EditField) GetField() Field {
	if x != nil && x.Field != nil {
		return *x.Field
	}
	return Field_UNDEFINED
}

func (x *EditField) GetReferrerMyPollsPage() int32 {
	if x != nil && x.ReferrerMyPollsPage != nil {
		return *x.ReferrerMyPollsPage
	}
	return 0
}

var File_internal_core_forecaster_telegram_proto_editfield_editfield_proto protoreflect.FileDescriptor

var file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDesc = []byte{
	0x0a, 0x41, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x67,
	0x72, 0x61, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x64, 0x69, 0x74, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x81,
	0x01, 0x0a, 0x09, 0x45, 0x64, 0x69, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x05, 0x52, 0x06, 0x70,
	0x6f, 0x6c, 0x6c, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x03,
	0x20, 0x02, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x33, 0x0a,
	0x16, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x5f, 0x6d, 0x79, 0x5f, 0x70, 0x6f, 0x6c,
	0x6c, 0x73, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x4d, 0x79, 0x50, 0x6f, 0x6c, 0x6c, 0x73, 0x50, 0x61,
	0x67, 0x65, 0x2a, 0x53, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x0d, 0x0a, 0x09, 0x55,
	0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x49,
	0x54, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x45, 0x53, 0x43, 0x52, 0x49, 0x50,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x52, 0x54, 0x5f,
	0x44, 0x41, 0x54, 0x45, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48,
	0x5f, 0x44, 0x41, 0x54, 0x45, 0x10, 0x04, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x65, 0x64, 0x69,
	0x74, 0x66, 0x69, 0x65, 0x6c, 0x64,
}

var (
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescOnce sync.Once
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescData = file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDesc
)

func file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescGZIP() []byte {
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescOnce.Do(func() {
		file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescData)
	})
	return file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDescData
}

var file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_goTypes = []interface{}{
	(Field)(0),        // 0: editfield.Field
	(*EditField)(nil), // 1: editfield.EditField
}
var file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_depIdxs = []int32{
	0, // 0: editfield.EditField.field:type_name -> editfield.Field
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_init() }
func file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_init() {
	if File_internal_core_forecaster_telegram_proto_editfield_editfield_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditField); i {
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
			RawDescriptor: file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_goTypes,
		DependencyIndexes: file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_depIdxs,
		EnumInfos:         file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_enumTypes,
		MessageInfos:      file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_msgTypes,
	}.Build()
	File_internal_core_forecaster_telegram_proto_editfield_editfield_proto = out.File
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_rawDesc = nil
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_goTypes = nil
	file_internal_core_forecaster_telegram_proto_editfield_editfield_proto_depIdxs = nil
}
