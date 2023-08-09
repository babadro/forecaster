package helpers

import (
	"fmt"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/golang/protobuf/proto"
)

func CallbackData(route byte, m proto.Message) (*string, error) {
	binaryData, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("can't marshal proto message: %w", err)
	}

	res := make([]byte, 0, len(binaryData)+1)
	res = append(res, route)
	res = append(res, binaryData...)

	return helpers.Ptr(string(res)), nil
}
