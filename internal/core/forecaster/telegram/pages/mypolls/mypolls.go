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
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
	w  dbwrapper.Wrapper
	allPolls  func(page int32) proto2.Message
	singlePoll func(itemID, _ int32) proto2.Message
}

func New(db models.DB) *Service {
	return &Service{
		db: db, w: dbwrapper.New(db),
		allPolls: func(page int32) proto2.Message {
			return &mypolls.MyPolls{CurrentPage: helpers.Ptr(page)}
		},
		singlePoll: func(itemID, _ int32) proto2.Message {
			return &editpoll.EditPoll{
				PollId: helpers.Ptr(itemID),
			}
		},
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

	keyboardIn := render.ManyItemsKeyboardInput{
		IDs:                       models2.PollsIDs(pollArr),
		CurrentPage:               currentPage,
		Prev:                      currentPage > 1,
		Next:                      currentPage*pageSize < totalCount,
		AllItemsRoute:             models.MyPollsRoute,
		SingleItemRoute:           models.EditPollRoute,
		AllItemsProtoMessage:      s.,
		SingleItemProtoMessage:    nil,
		FirstRowAdditionalButtons: nil,
	}

	keyboard, err := keyboardMarkup()
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(
		upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, txtMsg, keyboard), "", nil
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

func keyboardMarkup() (tgbotapi.InlineKeyboardMarkup, error) {
	editPollData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal editPoll callback data: %s", err.Error())
	}

	mainMenuButton, err := proto.MarshalCallbackData(models.MainPageRoute, &mainpage.MainPage{})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal mainPage callback data: %s", err.Error())
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "Create poll",
				CallbackData: editPollData,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "Main Menu",
				CallbackData: mainMenuButton,
			},
		),
	), nil
}
