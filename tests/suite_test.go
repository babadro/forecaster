package polls_tests

import (
	"context"
	"fmt"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/babadro/forecaster/tests/db"
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
	forecasterDB *postgres.ForecasterDB
	testDB       *db.TestDB

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

	s.forecasterDB = postgres.NewForecasterDB(dbPool)
	s.testDB = db.NewTestDB(dbPool)
}

func (s *APITestSuite) TearDownSuite() {

}

func (s *APITestSuite) SetupTest() {

}

func (s *APITestSuite) TearDownTest() {

}

func (s *APITestSuite) CreateDefaultSeries() {
	s.T().Helper()
	series := models.Series{
		ID:          0,
		Description: "default series desc",
		Title:       "default series title",
	}

	_, err := s.testDB.CreateSeries(context.Background(), series)
	s.Require().NoError(err)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
