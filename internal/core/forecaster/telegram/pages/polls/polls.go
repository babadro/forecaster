package polls

import (
	"context"
	"fmt"
	"strconv"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	models3 "github.com/babadro/forecaster/internal/models"

	models2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/helpers"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db         models.DB
	allPolls   func(page int32) proto2.Message
	singlePoll func(itemID, _ int32) proto2.Message
}

func New(db models.DB) *Service {
	return &Service{
		db: db,
		allPolls: func(page int32) proto2.Message {
			return &polls.Polls{CurrentPage: helpers.Ptr(page)}
		},
		singlePoll: func(itemID, _ int32) proto2.Message {
			return &poll.Poll{
				PollId: helpers.Ptr(itemID),
			}
		},
	}
}

func (s *Service) NewRequest() (proto2.Message, *polls.Polls) {
	v := new(polls.Polls)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	currentPageStr := upd.Message.Text[len(models.ShowPollsStartCommandPrefix):]

	currentPage, err := strconv.ParseInt(currentPageStr, 10, 32)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse current page: %s", err.Error())
	}

	return s.render(ctx, int32(currentPage), upd.Message.Chat.ID, upd.Message.MessageID, false, false)
}

func (s *Service) RenderCallback(
	ctx context.Context, req *polls.Polls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	filterFinished := models.PollsFlags(req.GetFlags()).IsSet(models.FilterFinishedStatus)

	return s.render(ctx, req.GetCurrentPage(), chat.ID, message.MessageID, true, filterFinished)
}

const pageSize = 10

func (s *Service) render(
	ctx context.Context, currentPage int32, chatID int64, messageID int, editMessage bool, filterFinished bool,
) (tgbotapi.Chattable, string, error) {
	offset, limit := uint64((currentPage-1)*pageSize), uint64(pageSize)

	status := models3.ActivePollStatus
	if filterFinished {
		status = models3.FinishedPollStatus
	}

	pollsArr, totalCount, err := s.db.GetPolls(ctx, offset, limit,
		models.NewPollFilter().WithStatus(status),
		models.PollSort{
			By:  models.PopularityPollSort,
			Asc: false,
		})
	if err != nil {
		return nil, "", fmt.Errorf("unable to get polls: %s", err.Error())
	}

	filterBtn, err := filterButton(!filterFinished)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create filter button: %s", err.Error())
	}

	keyboardIn := render.ManyItemsKeyboardInput{
		IDs:                       models2.PollsIDs(pollsArr),
		CurrentPage:               currentPage,
		Prev:                      currentPage > 1,
		Next:                      currentPage*pageSize < totalCount,
		AllItemsRoute:             models.PollsRoute,
		SingleItemRoute:           models.PollRoute,
		AllItemsProtoMessage:      s.allPolls,
		SingleItemProtoMessage:    s.singlePoll,
		FirstRowAdditionalButtons: []tgbotapi.InlineKeyboardButton{filterBtn},
	}

	keyboard, err := render.ManyItemsKeyboardMarkup(keyboardIn)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txtMsg(pollsArr), keyboard),
			"", nil
	}

	return render.NewMessageWithKeyboard(chatID, txtMsg(pollsArr), keyboard), "", nil
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

func filterButton(filterFinished bool) (tgbotapi.InlineKeyboardButton, error) {
	btnText, flags := "Show active", int32(0)
	if filterFinished {
		btnText, flags = "Show finished", int32(models.FilterFinishedStatus)
	}

	data, err := proto.MarshalCallbackData(models.PollsRoute, &polls.Polls{
		CurrentPage: helpers.Ptr[int32](1),
		Flags:       helpers.Ptr(flags),
	})
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("unable to marshal callback data: %s", err.Error())
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         btnText,
		CallbackData: data,
	}, nil
}
