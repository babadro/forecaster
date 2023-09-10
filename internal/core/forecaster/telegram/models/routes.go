package models

const (
	MainPageRoute byte = iota
	VotePreviewRoute
	VoteRoute
	PollRoute
	UserPollResultRoute
)

const (
	ShowPollStartCommandPrefix  = "/start showpoll_"
	ShowUserResultCommandPrefix = "/start showuserres_"
)
