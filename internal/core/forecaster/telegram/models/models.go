package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MaxCountInRow = 8

	Days365     = 365
	Hours24     = 24
	Seconds3600 = 3600

	Percent100 = 100
)

type TgBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type EditButton[T any] struct {
	Text  string
	Field T
}
