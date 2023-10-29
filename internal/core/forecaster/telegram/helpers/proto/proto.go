package proto

import (
	"fmt"

	"encoding/base64"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editstatus"
	"github.com/babadro/forecaster/internal/helpers"
	models2 "github.com/babadro/forecaster/internal/models"
	"google.golang.org/protobuf/proto"
)

const minCallbackDataLength = 1

func MarshalCallbackData(route byte, m proto.Message) (*string, error) {
	binaryData, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("can't marshal proto message: %w", err)
	}

	res := make([]byte, 0, len(binaryData)+1)
	res = append(res, route)
	res = append(res, binaryData...)

	encoded := base64.StdEncoding.EncodeToString(res)

	return helpers.Ptr(encoded), nil
}

func UnmarshalCallbackData(data string, m proto.Message) error {
	if len(data) < minCallbackDataLength {
		return fmt.Errorf("callback data is too short")
	}

	route := data[0]

	if err := proto.Unmarshal([]byte(data[1:]), m); err != nil {
		return fmt.Errorf("can't unmarshal proto message for route %d: %s", route, err.Error())
	}

	return nil
}

func PollStatusFromProto(in editstatus.Status) (models2.PollStatus, error) {
	switch in {
	case editstatus.Status_UNKNOWN:
		return models2.UnknownPollStatus, nil
	case editstatus.Status_DRAFT:
		return models2.DraftPollStatus, nil
	case editstatus.Status_ACTIVE:
		return models2.ActivePollStatus, nil
	case editstatus.Status_FINISHED:
		return models2.FinishedPollStatus, nil
	default:
		return models2.UnknownPollStatus, fmt.Errorf("unknown poll status %d", in)
	}
}
