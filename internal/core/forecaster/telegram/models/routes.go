package models

const (
	MainPageRoute byte = iota
	VotePreviewRoute
	VoteRoute
	PollRoute
	UserPollResultRoute
	PollsRoute
)

const (
	ShowPollStartCommandPrefix  = "/start showpoll_"
	ShowUserResultCommandPrefix = "/start showuserres_"
	ShowPollsStartCommandPrefix = "/start showpolls_"
	ShowForecastsStartCommand   = "/start showforecasts"
)
