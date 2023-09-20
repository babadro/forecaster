package telegram_test

import (
	"regexp"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Vote 2 times for different options of the same poll and verify the results...
func (s *TelegramServiceSuite) TestVoting() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll(time.Now())

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID, 456)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	// each keyboard button is a poll option
	s.Require().Len(pollButtons, len(poll.Options)+1) // +1 for "All Polls" button

	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton, 456)

	// verify the result votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)

	// verify contains the poll title and description
	txt := votePreviewMsg.Text
	option := s.findOptionByCallbackData(poll, firstButton.CallbackData)
	s.Require().Contains(txt, option.Title)
	s.Require().Contains(txt, option.Description)

	// verify message has two buttons
	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0], 456)

	// verify the vote message
	voteMsg := s.asEditMessage(sentMsg)

	s.Require().Contains(voteMsg.Text, "Success")

	// push back to poll button
	voteKeyboard := getButtons(*voteMsg.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton := voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton, 456)

	// verify the poll message
	pollMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg2.Text, poll.Title)
	pattern := "Last time you voted for:.+" + option.Title
	regex := regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(pollMsg2.Text), "expected %s to match regex %s", pollMsg2.Text, pattern)

	// each keyboard button is a poll option
	pollButtons2 := getButtons(*pollMsg2.ReplyMarkup)
	s.Require().Len(pollButtons2, len(poll.Options)+1) // +1 for "All Polls" button

	// chose option I didn't vote earlier
	anotherOptionButton, found := tgbotapi.InlineKeyboardButton{}, false

	for _, button := range pollButtons2 {
		op := s.findOptionByCallbackData(poll, button.CallbackData)
		if op.ID != option.ID {
			anotherOptionButton, found = button, true
			break
		}
	}

	s.Require().True(found)

	// sleep for second to make sure vote timestamp (which used second precision) is different
	time.Sleep(time.Second)
	// push the button to vote for another option this time
	s.sendCallback(anotherOptionButton, 456)

	// verify the votepreview message
	votePreviewMsg2 := s.asEditMessage(sentMsg)

	// verify the poll contains title and description
	txt = votePreviewMsg2.Text
	anotherOption := s.findOptionByCallbackData(poll, anotherOptionButton.CallbackData)
	s.Require().Contains(txt, anotherOption.Title)
	s.Require().Contains(txt, anotherOption.Description)

	// verify message has two buttons
	votePreviewButtons = getButtons(*votePreviewMsg2.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0], 456)

	// verify the vote message
	voteMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(voteMsg2.Text, "Success")

	// push back to poll button
	voteKeyboard = getButtons(*voteMsg2.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton = voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton, 456)

	// verify the poll message
	pollMsg3 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg3.Text, poll.Title)
	pattern = "Last time you voted for:.+" + anotherOption.Title
	regex = regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(pollMsg3.Text), "expected %s to match regex: %q", pollMsg3.Text, pattern)

	// each keyboard button is a poll option
	pollButtons3 := getButtons(*pollMsg3.ReplyMarkup)
	s.Require().Len(pollButtons3, len(poll.Options)+1) // +1 for "All Polls" button
}

func (s *TelegramServiceSuite) TestVotePreview_BackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll(time.Now())

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID, 456)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)

	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton, 456)

	// verify the result votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)

	// verify message has two buttons
	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the back button
	s.sendCallback(votePreviewButtons[1], 456)

	// verify the poll message
	pollMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg2.Text, poll.Title)
}

func (s *TelegramServiceSuite) Test_expiredPoll() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0
	pollInput.Finish = strfmt.DateTime(time.Now().Add(-time.Hour)) // expired

	poll := s.createPoll(pollInput, time.Now())

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID, 456)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)
	// verify the poll message
	s.Require().Contains(pollMsg.Text, "poll has expired")

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	// send the first option
	s.sendCallback(pollButtons[0], 456)

	// verify votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(votePreviewMsg.Text, "poll is expired")

	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 1)
	// the only button is "Back"
	s.Require().Contains(votePreviewButtons[0].Text, "Back")
}

func (s *TelegramServiceSuite) Test_attempt_to_vote_for_the_same_option_result_in_error() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll(time.Now())

	userID := int64(456)

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID, userID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton, userID)

	// verify the result votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)
	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().NotEmpty(votePreviewButtons)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0], userID)

	// verify the vote message
	voteMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(voteMsg.Text, "Success")

	// try to vote again
	s.sendCallback(votePreviewButtons[0], userID)

	// verify the vote message
	voteMsg2 := s.asEditMessage(sentMsg)
	s.Require().Contains(voteMsg2.Text, "You have already voted for this option")

	// verify keyboard
	voteKeyboard := getButtons(*voteMsg2.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)
	backButton := voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
}
