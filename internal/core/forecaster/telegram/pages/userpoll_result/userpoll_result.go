package userpollresult

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db      models.DB
	w       dbwrapper.Wrapper
	botName string
}

func New(db models.DB, botName string) *Service {
	return &Service{db: db, w: dbwrapper.New(db), botName: botName}
}

func (s *Service) NewRequest() (proto2.Message, *userpollresult.UserPollResult) {
	v := new(userpollresult.UserPollResult)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	pollID, userID, err := parseIDs(upd.Message.Text[len(models.ShowUserResultCommandPrefix):])
	if err != nil {
		return nil, "", err
	}

	return s.render(
		ctx, pollID, userID, upd.Message.Chat.ID, upd.Message.MessageID, upd.Message.From.UserName, false,
	)
}

func parseIDs(text string) (int32, int64, error) {
	ids := strings.Split(text, "_")
	if length := len(ids); length != 2 {
		return 0, 0, fmt.Errorf(
			"userpoll result: can't parse pollID and userID from %s, expected len(ids)=2, got=%d", text, length)
	}

	pollID, err := strconv.ParseInt(ids[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("userpoll result: can't parse pollID from command: %s", text)
	}

	userID, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("userpoll result: can't parse userID from command: %s", text)
	}

	return int32(pollID), userID, nil
}

func (s *Service) RenderCallback(
	ctx context.Context, req *userpollresult.UserPollResult, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	user := upd.CallbackQuery.From
	if user == nil {
		return nil, "", errors.New("user is nil")
	}

	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	return s.render(ctx, req.GetPollId(), user.ID, chat.ID, message.MessageID, user.UserName, true)
}

func (s *Service) render(
	ctx context.Context,
	pollID int32,
	userID, chatID int64,
	messageID int,
	userName string,
	editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	outcome, idx := swagger.GetOutcome(p.Options)
	if idx == -1 {
		return nil, "", fmt.Errorf("userpoll result: can't get outcome for pollID: %d", p.ID)
	}

	userVote, found, err := s.w.GetUserVote(ctx, userID, p.ID)
	if err != nil {
		return nil, "", err
	}

	if !found {
		return nil, "", fmt.Errorf("userpoll result: can't find last user's vote for pollID: %d", p.ID)
	}

	if userVote.OptionID != outcome.ID {
		return nil, "", fmt.Errorf("userpoll result: last user's vote is not outcome for pollID: %d", p.ID)
	}

	stat := getStatistic(p.Options, userVote.Position)

	txtInputModel := txtMsgInput{
		userName:                userName,
		userID:                  userID,
		pollID:                  p.ID,
		optionTitle:             outcome.Title,
		finishPoll:              time.Time(p.Finish),
		voteUnixTime:            userVote.EpochUnixTimestamp,
		prozentOfAllVotesBehind: stat.prozentOfAllVotesBehind,
		prozentOfWonVotesBehind: stat.prozentOfWonVotesBehind,
		totalVotes:              stat.totalVotes,
		totalVotesForWonOption:  stat.votesForWonOption,
	}

	msg := s.txtMsg(txtInputModel)

	markup, err := keyboardMarkup(p.ID)
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to create keyboard markup: %s", err.Error())
	}

	var res tgbotapi.Chattable
	if editMessage {
		res = render.NewEditMessageTextWithKeyboard(chatID, messageID, msg, markup)
	} else {
		res = render.NewMessageWithKeyboard(chatID, msg, markup)
	}

	return res, "", nil
}

type statistics struct {
	prozentOfAllVotesBehind int8
	prozentOfWonVotesBehind int8
	votesForWonOption       int32
	totalVotes              int32
}

func getStatistic(options []*swagger.Option, userPositionAmongWonVotes int32) statistics {
	var votesForLoseOptions, votesForWonOption int32

	for _, o := range options {
		if !o.IsActualOutcome {
			votesForLoseOptions += o.TotalVotes
		} else {
			votesForWonOption = o.TotalVotes
		}
	}

	totalVotes := votesForLoseOptions + votesForWonOption

	numberOfVotesBehind := votesForLoseOptions + votesForWonOption - userPositionAmongWonVotes

	prozentOfAllVotesBehind := int8(float32(numberOfVotesBehind) / float32(totalVotes) * 100)

	prozentOfWonVotesBehind := int8(float32(votesForWonOption-userPositionAmongWonVotes) / float32(votesForWonOption) * 100)

	return statistics{
		prozentOfAllVotesBehind: prozentOfAllVotesBehind,
		prozentOfWonVotesBehind: prozentOfWonVotesBehind,
		votesForWonOption:       votesForWonOption,
		totalVotes:              totalVotes,
	}
}

type txtMsgInput struct {
	userName                string
	userID                  int64
	pollID                  int32
	optionTitle             string
	finishPoll              time.Time
	voteUnixTime            int64
	prozentOfAllVotesBehind int8
	prozentOfWonVotesBehind int8
	totalVotes              int32
	totalVotesForWonOption  int32
}

func (s *Service) txtMsg(in txtMsgInput) string {
	var sb render.StringBuilder

	advanceTimeNumber, advanceTimeUnit := render.GetHighestTimeUnit(in.finishPoll.Sub(time.Unix(in.voteUnixTime, 0)))

	sb.Printf("<b>%s</b> you predicted that %s %d %s before!", in.userName, in.optionTitle, advanceTimeNumber, advanceTimeUnit)

	if in.prozentOfAllVotesBehind != 0 && in.prozentOfWonVotesBehind != 0 {
		sb.Printf("\nThis places you ahead of %d%% of all participants and shows that you chose the correct option earlier than %d%% of those who also chose correctly.", in.prozentOfAllVotesBehind, in.prozentOfWonVotesBehind)
	} else if in.prozentOfAllVotesBehind != 0 {
		sb.Printf("\nThis places you ahead of %d%% of all participants", in.prozentOfAllVotesBehind)
	}

	sb.Printf("\nOut of %d total participants, only %d made a correct prediction.", in.totalVotes, in.totalVotesForWonOption)

	sb.Print("\nShare your results by forwarding this message or by sending the following link:")
	sb.Printf("https://t.me/%s?start=show_user_result_%d_%d", s.botName, in.pollID, in.userID)

	return sb.String()
}

func keyboardMarkup(pollID int32) (tgbotapi.InlineKeyboardMarkup, error) {
	backData, err := proto.MarshalCallbackData(models.PollRoute, &poll.Poll{PollId: helpers.Ptr[int32](pollID)})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall poll callback data: %s", err.Error())
	}

	backBtn := tgbotapi.InlineKeyboardButton{Text: "Back", CallbackData: backData}

	return render.Keyboard(backBtn), nil
}
