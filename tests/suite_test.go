package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"testing"
	"time"

	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/caarlos0/env"
	"github.com/stretchr/testify/suite"
)

type envVars struct {
	AppPort int    `env:"APP_PORT,required"`
	DBConn  string `env:"DB_CONN,required"`
}

var (
	envs    envVars
	apiAddr string
)

// APITestSuite defines the suite
type APITestSuite struct {
	suite.Suite
	db *postgres.ForecasterDB

	client *http.Client
}

// SetupSuite function will be run by testify before any tests or test suites are run.
func (s *APITestSuite) SetupSuite() {
	s.Require().NoError(env.Parse(&envs))

	apiAddr = fmt.Sprintf("http://localhost:%d", envs.AppPort)

	s.client = &http.Client{
		Timeout: time.Second * 10,
	}

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}

	s.db = postgres.NewForecasterDB(dbPool)

}

func (s *APITestSuite) TearDownSuite() {

}

func (s *APITestSuite) SetupTest() {

}

func (s *APITestSuite) TearDownTest() {

}

func (s *APITestSuite) CreateDefaultSeries() {
	s.db.CreatePoll()

	series := models.CreateSeries{
		Description: "test desc",
		Title:       "test title",
	}

	b, err := json.Marshal(series)
	s.Require().NoError(err)

	body := bytes.NewReader(b)

	resp, err := s.client.Post(apiAddr+"/series", "application/json", body)
	s.Require().NoError(err)

}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
