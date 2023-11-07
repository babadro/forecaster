package models

type ForecastsFlags int32

const (
	FilterFinishedForecasts ForecastsFlags = 1 << iota
)

func (f ForecastsFlags) IsSet(flag ForecastsFlags) bool {
	return f&flag != 0
}
