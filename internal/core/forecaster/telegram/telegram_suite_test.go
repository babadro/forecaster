package telegram_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram"
	"github.com/babadro/forecaster/internal/infra/postgres"
	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/babadro/forecaster/mocks"
	"github.com/babadro/forecaster/tests/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type envVars struct {
	DBConn    string `env:"DB_CONN" envDefault:"postgres://postgres:postgres@localhost:5432/forecaster?sslmode=disable"`
	SleepMode bool   `env:"SLEEP_MODE" envDefault:"false"`
}

type TelegramServiceSuite struct {
	suite.Suite
	db        *postgres.ForecasterDB
	testDB    *db.TestDB
	mockTgBot *mocks.TelegramBot

	telegramService *telegram.Service
	logOutput       bytes.Buffer
	logger          zerolog.Logger

	sleepMode bool
}

func (s *TelegramServiceSuite) SetupSuite() {
	var envs envVars

	s.Require().NoError(env.Parse(&envs))

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	s.Require().NoError(err)

	s.testDB = db.NewTestDB(dbPool)
	s.db = postgres.NewForecasterDB(dbPool)

	s.sleepMode = envs.SleepMode
}

func (s *TelegramServiceSuite) SetupTest() {
	s.createDefaultSeries()

	s.mockTgBot = &mocks.TelegramBot{}

	s.logOutput = bytes.Buffer{}
	s.logger = zerolog.New(&s.logOutput)

	s.telegramService = telegram.NewService(s.db, s.mockTgBot, "test-bot")
}

func (s *TelegramServiceSuite) TearDownTest() {
	if s.sleepMode {
		time.Sleep(time.Hour * 10_000)
	}

	s.cleanAllTables()
}

func (s *TelegramServiceSuite) cleanAllTables() {
	s.T().Helper()

	s.Require().NoError(s.testDB.CleanAllTables(context.Background()))
}

func (s *TelegramServiceSuite) createDefaultSeries() {
	s.T().Helper()

	series := models.Series{
		ID:          0,
		Description: "default series desc",
		Title:       "default series title",
	}

	_, err := s.testDB.CreateSeries(context.Background(), series)
	s.Require().NoError(err)
}

func randomModel[T any](t *testing.T) T {
	var model T

	require.NoError(t, gofakeit.Struct(&model))

	return model
}

func TestTelegram(t *testing.T) {
	suite.Run(t, new(TelegramServiceSuite))
}
