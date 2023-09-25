package polls

import (
	"context"
	"fmt"
	"strconv"

	"github.com/babadro/forecaster/internal/helpers"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
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

	return s.render(ctx, int32(currentPage), upd.Message.Chat.ID, upd.Message.MessageID, false)
}

func (s *Service) RenderCallback(
	ctx context.Context, req *polls.Polls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	return s.render(ctx, req.GetCurrentPage(), chat.ID, message.MessageID, true)
}

const pageSize = 10

func (s *Service) render(
	ctx context.Context, currentPage int32, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	offset, limit := uint64((currentPage-1)*pageSize), uint64(pageSize)

	pollsArr, totalCount, err := s.db.GetPolls(ctx, offset, limit)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get polls: %s", err.Error())
	}

	keyboardIn := render.KeyboardInput{
		IDs:         pollsIDs(pollsArr),
		CurrentPage: currentPage,
		Prev:        currentPage > 1,
		Next:        currentPage*pageSize < totalCount,
		Route:       models.PollRoute,
		AllItemsProtoMessage: func(page int32) proto2.Message {
			return &polls.Polls{
				CurrentPage: helpers.Ptr(page),
			}
		},
	}

	keyboard, err := render.KeyboardMarkup(keyboardIn)
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

func pollsIDs(p []swagger.Poll) []int32 {
	ids := make([]int32, len(p))
	for i, poll := range p {
		ids[i] = poll.ID
	}

	return ids
}
