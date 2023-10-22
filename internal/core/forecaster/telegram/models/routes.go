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
	EditPollRoute
	MyPollsRoute
	EditPollFieldRoute
	EditOptionRoute
	EditOptionFieldRoute
	DeletePollRoute
	DeleteOptionRoute
)

const (
	ShowMainStartCommandPrefix      = "/start main"
	ShowPollStartCommandPrefix      = "/start showpoll_"
	ShowUserResultCommandPrefix     = "/start showuserres_"
	ShowPollsStartCommandPrefix     = "/start showpolls_"
	ShowForecastsStartCommandPrefix = "/start showforecasts_"
	ShowForecastStartCommandPrefix  = "/start showforecast_"
)

const (
	EditPollCommand   = "/editpoll"
	EditOptionCommand = "/editoption"
)
