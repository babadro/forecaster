package mypolls

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	models2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
	w  dbwrapper.Wrapper
}

func New(db models.DB) *Service {
	return &Service{
		db: db, w: dbwrapper.New(db),
	}
}

func (s *Service) NewRequest() (proto2.Message, *mypolls.MyPolls) {
	v := new(mypolls.MyPolls)

	return v, v
}

const pageSize = 10

func (s *Service) RenderCallback(
	ctx context.Context, req *mypolls.MyPolls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	currentPage := req.GetCurrentPage()
	if currentPage == 0 {
		currentPage = 1
	}

	offset, limit := uint64((currentPage-1)*pageSize), uint64(pageSize)

	userID := upd.CallbackQuery.From.ID

	pollArr, totalCount, err := s.db.GetPolls(ctx, offset, limit, models.NewPollFilter().WithTelegramUserID(userID))
	if err != nil {
		return nil, "", fmt.Errorf("unable to get polls: %s", err.Error())
	}

	additionalBtns, err := additionalButtons()
	if err != nil {
		return nil, "", fmt.Errorf("unable to create additional buttons: %s", err.Error())
	}

	keyboardIn := render.ManyItemsKeyboardInput{
		IDs:             models2.PollsIDs(pollArr),
		CurrentPage:     currentPage,
		Prev:            currentPage > 1,
		Next:            currentPage*pageSize < totalCount,
		AllItemsRoute:   models.MyPollsRoute,
		SingleItemRoute: models.EditPollRoute,
		AllItemsProtoMessage: func(page int32) proto2.Message {
			return &mypolls.MyPolls{CurrentPage: helpers.Ptr(page)}
		},
		SingleItemProtoMessage: func(itemID, referrerMyPollsPage int32) proto2.Message {
			return &editpoll.EditPoll{
				PollId:              helpers.Ptr(itemID),
				ReferrerMyPollsPage: helpers.Ptr(referrerMyPollsPage),
			}
		},
		FirstRowAdditionalButtons: additionalBtns,
	}

	keyboard, err := render.ManyItemsKeyboardMarkup(keyboardIn)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(
			upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, txtMsg(pollArr), keyboard),
		"", nil
}

func txtMsg(pollsArr []swagger.Poll) string {
	if len(pollsArr) == 0 {
		return "There are no polls yet"
	}

	var sb render.StringBuilder
	for i, p := range pollsArr {
		sb.Printf("%d. %s\n", i+1, p.Title)
	}

	return sb.String()
}

func additionalButtons() ([]tgbotapi.InlineKeyboardButton, error) {
	editPollData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal editPoll callback data: %s", err.Error())
	}

	return []tgbotapi.InlineKeyboardButton{
		{
			Text:         "Create poll",
			CallbackData: editPollData,
		},
	}, nil
}
