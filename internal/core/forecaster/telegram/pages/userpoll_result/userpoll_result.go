package userpollresult

import (
	"context"
	"fmt"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
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

func (s *Service) RenderCallback(
	ctx context.Context, req *userpollresult.UserPollResult, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, req.GetPollId())
	if err != nil {
		return nil, errMsg, err
	}

	outcome, idx := swagger.GetOutcome(p.Options)
	if idx == -1 {
		return nil, "", fmt.Errorf("userpoll result: can't get outcome for pollID: %d", p.ID)
	}

	user := upd.CallbackQuery.From

	userVote, found, err := s.w.GetUserVote(ctx, user.ID, p.ID)
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
		userName:                user.UserName,
		userID:                  user.ID,
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

	markup, err := keyboardMarkup()
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to create keyboard markup: %s", err.Error())
	}

	origMsg := upd.CallbackQuery.Message

	return render.NewEditMessageTextWithKeyboard(origMsg.Chat.ID, origMsg.MessageID, msg, markup), "", nil
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

	sb.Printf("\nThis places you ahead of %d%% of all participants and shows that you chose the correct option earlier than %d%% of those who also chose correctly.", in.prozentOfAllVotesBehind, in.prozentOfWonVotesBehind)

	sb.Printf("\nOut of %d total participants, only %d made a correct prediction.", in.totalVotes, in.totalVotesForWonOption)

	sb.Print("\nShare your results by forwarding this message or by sending the following link:")
	sb.Printf("https://t.me/%s?start=show_user_result_%d_%d", s.botName, in.pollID, in.userID)

	return sb.String()
}

func keyboardMarkup() (tgbotapi.InlineKeyboardMarkup, error) {
	return render.Keyboard(
		tgbotapi.InlineKeyboardButton{
			// todo
		},
	), nil
}
