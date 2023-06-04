package forecaster_bot

import bot "github.com/babadro/forecaster/pkg/forecaster_bot"

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db: db}
}

type DB interface {
	GetByID(id int) (bot.Poll, error)
}

func (s *Service) GetByID(id int) (bot.Poll, error) {
	return s.db.GetByID(id)
}
