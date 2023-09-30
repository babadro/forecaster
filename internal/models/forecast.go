package models

type Forecast struct {
	PollID    int32
	PollTitle string
	Options   []ForecastOption
}

type ForecastOption struct {
	ID         int16
	Title      string
	TotalVotes int32
}
