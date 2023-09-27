package models

const (
	MainPageRoute byte = iota
	VotePreviewRoute
	VoteRoute
	PollRoute
	UserPollResultRoute
	PollsRoute
	ForecastsRoute
	ForecastRoute
)

const (
	ShowPollStartCommandPrefix      = "/start showpoll_"
	ShowUserResultCommandPrefix     = "/start showuserres_"
	ShowPollsStartCommandPrefix     = "/start showpolls_"
	ShowForecastsStartCommandPrefix = "/start showforecasts_"
	ShowForecastStartCommandPrefix  = "/start showforecast_"
)
