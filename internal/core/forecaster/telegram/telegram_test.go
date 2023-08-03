package telegram_test

import (
	"context"

	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type envVars struct {
	DBConn string `env:"DB_CONN,required"`
}

type TelegramServiceSuite struct {
	suite.Suite
	db     *postgres.ForecasterDB
	dbPool *pgxpool.Pool
}

func (s *TelegramServiceSuite) SetupSuite() {
	var envs envVars

	s.Require().NoError(env.Parse(&envs))

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	s.Require().NoError(err)

	s.dbPool = dbPool
	s.db = postgres.NewForecasterDB(dbPool)

}

func (s *TelegramServiceSuite) TearDownTest() {
	s.cleanAllTables()
}

func (s *TelegramServiceSuite) cleanAllTables() {
	s.T().Helper()

	for _, tableName := range []string{
		"forecaster.series",
		"forecaster.polls",
		"forecaster.options",
	} {
		_, err := s.dbPool.Exec(context.Background(), "TRUNCATE TABLE "+tableName+" CASCADE")
		s.Require().NoError(err)
	}
}
