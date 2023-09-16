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

const idsCount = 2

func parseIDs(text string) (int32, int64, error) {
	ids := strings.Split(text, "_")
	if length := len(ids); length != idsCount {
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
	isCallback bool,
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
		return nil, "", fmt.Errorf("userpoll result: can't find user's vote for pollID: %d", p.ID)
	}

	markup, err := keyboardMarkup(p.ID)
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to create keyboard markup: %s", err.Error())
	}

	if userVote.OptionID != outcome.ID {
		if !isCallback { // it is assumed that unsuccessful vote is not possible to share via start command
			return nil, "", fmt.Errorf("userpoll result: user's vote is not outcome for pollID: %d", p.ID)
		}

		votesForLoseOptions, votesForWonOption := getGeneralStatistic(p.Options)

		return render.NewEditMessageTextWithKeyboard(
			chatID, messageID,
			s.txtMsgForWrongVotedUser(userName, outcome.Title, votesForLoseOptions+votesForWonOption, votesForWonOption),
			markup), "", nil
	}

	stat, err := getUserStatistic(p.Options, userVote.Position)
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to get user statistic: %s", err.Error())
	}

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
		votesForWonOption:       stat.votesForWonOption,
	}

	var res tgbotapi.Chattable

	if isCallback {
		msg := s.txtMsg(txtInputModel)
		res = render.NewEditMessageTextWithKeyboard(chatID, messageID, msg, markup)
	} else {
		msg := s.thirdPersonTxtMsg(txtInputModel)
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

const prozents100 = 100

func getUserStatistic(options []*swagger.Option, userPositionAmongWonVotes int32) (statistics, error) {
	votesForLoseOptions, votesForWonOption := getGeneralStatistic(options)

	totalVotes := votesForLoseOptions + votesForWonOption

	if totalVotes == 0 {
		return statistics{}, errors.New("total votes is 0")
	}

	if votesForWonOption == 0 {
		return statistics{}, errors.New("votes for won option is 0")
	}

	numberOfVotesBehind := votesForLoseOptions + votesForWonOption - userPositionAmongWonVotes

	prozentOfAllVotesBehind := int8(float32(numberOfVotesBehind) / float32(totalVotes) * prozents100)

	prozentOfWonVotesBehind :=
		int8(float32(votesForWonOption-userPositionAmongWonVotes) / float32(votesForWonOption) * prozents100)

	return statistics{
		prozentOfAllVotesBehind: prozentOfAllVotesBehind,
		prozentOfWonVotesBehind: prozentOfWonVotesBehind,
		votesForWonOption:       votesForWonOption,
		totalVotes:              totalVotes,
	}, nil
}

func getGeneralStatistic(options []*swagger.Option) (int32, int32) {
	var votesForLoseOptions, votesForWonOption int32

	for _, o := range options {
		if !o.IsActualOutcome {
			votesForLoseOptions += o.TotalVotes
		} else {
			votesForWonOption = o.TotalVotes
		}
	}

	return votesForLoseOptions, votesForWonOption
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
	votesForWonOption       int32
}

const showUserResultCommand = "showuserres_"

func (s *Service) txtMsg(in txtMsgInput) string {
	var sb render.StringBuilder

	advanceTimeNumber, advanceTimeUnit := render.GetHighestTimeUnit(in.finishPoll.Sub(time.Unix(in.voteUnixTime, 0)))

	sb.Printf("<b>%s</b> you predicted that %s %d %s before!",
		in.userName, in.optionTitle, advanceTimeNumber, advanceTimeUnit)

	if in.prozentOfAllVotesBehind != 0 && in.prozentOfWonVotesBehind != 0 {
		sb.Printf(
			"\nThis places you ahead of %d%% of all participants "+
				"and shows that you chose the correct option earlier than %d%% of those who also chose correctly.",
			in.prozentOfAllVotesBehind, in.prozentOfWonVotesBehind)
	} else if in.prozentOfAllVotesBehind != 0 {
		sb.Printf("\nThis places you ahead of %d%% of all participants", in.prozentOfAllVotesBehind)
	}

	if in.votesForWonOption > 0 {
		sb.Printf("\nOut of %d total participants, only %d made a correct prediction.", in.totalVotes, in.votesForWonOption)
	} else {
		sb.Printf("\nOut of %d total participants, no one made a correct prediction.", in.totalVotes)
	}

	sb.Print("\nShare your results by forwarding this message or by sending the following link:")
	sb.Printf("https://t.me/%s?start=%s%d_%d", s.botName, showUserResultCommand, in.pollID, in.userID)

	return sb.String()
}

func (s *Service) txtMsgForWrongVotedUser(
	userName, votedOptionTitle string, totalVotes, votesForWonOption int32) string {
	var sb render.StringBuilder

	sb.Printf("<b>%s</b>, your prediction for '%s' didn't quite pan out this time.", userName, votedOptionTitle)

	if votesForWonOption > 0 {
		sb.Printf("\nOut of %d total participants, only %d made a correct prediction.", totalVotes, votesForWonOption)
	} else {
		sb.Printf("\nOut of %d total participants, no one made a correct prediction.", totalVotes)
	}

	return sb.String()
}

func (s *Service) thirdPersonTxtMsg(in txtMsgInput) string {
	var sb render.StringBuilder

	advanceTimeNumber, advanceTimeUnit := render.GetHighestTimeUnit(in.finishPoll.Sub(time.Unix(in.voteUnixTime, 0)))

	sb.Printf("<b>%s</b> predicted that %s %d %s before!", in.userName, in.optionTitle, advanceTimeNumber, advanceTimeUnit)

	if in.prozentOfAllVotesBehind != 0 && in.prozentOfWonVotesBehind != 0 {
		sb.Printf(
			"\nThis places them ahead of %d%% of all participants"+
				"and shows that %s chose the correct option earlier than %d%% of those who also chose correctly.",
			in.prozentOfAllVotesBehind, in.userName, in.prozentOfWonVotesBehind)
	} else if in.prozentOfAllVotesBehind != 0 {
		sb.Printf("\nThis places them ahead of %d%% of all participants", in.prozentOfAllVotesBehind)
	}

	if in.votesForWonOption > 0 {
		sb.Printf("\nOut of %d total participants, only %d made a correct prediction.", in.totalVotes, in.votesForWonOption)
	} else {
		sb.Printf("\nOut of %d total participants, no one made a correct prediction.", in.totalVotes)
	}

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
