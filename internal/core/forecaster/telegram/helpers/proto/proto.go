package proto

import (
	"fmt"

	"github.com/babadro/forecaster/internal/helpers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const minCallbackDataLength = 2

func MarshalCallbackData(route byte, m proto.Message) (*string, error) {
	binaryData, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("can't marshal proto message: %w", err)
	}

	res := make([]byte, 0, len(binaryData)+1)
	res = append(res, route)
	res = append(res, binaryData...)

	return helpers.Ptr(string(res)), nil
}

func UnmarshalCallbackData(data string, m proto.Message) error {
	if len(data) < minCallbackDataLength {
		return fmt.Errorf("callback data is too short")
	}

	route := data[0]

	binaryData := []byte(data[1:])

	if err := proto.Unmarshal(binaryData, m); err != nil {
		return fmt.Errorf("can't unmarshal proto message for route %d: %s", route, err.Error())
	}

	return nil
}

type Req interface {
	ProtoReflect() protoreflect.Message
}

func UnmarshalCallbackData2[T proto.Message](data string, m T) error {
	if len(data) < minCallbackDataLength {
		return fmt.Errorf("callback data is too short")
	}

	route := data[0]

	binaryData := []byte(data[1:])

	if err := proto.Unmarshal(binaryData, m); err != nil {
		return fmt.Errorf("can't unmarshal proto message for route %d: %s", route, err.Error())
	}

	return nil
}
