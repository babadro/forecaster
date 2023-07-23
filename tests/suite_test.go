package polls_tests

import (
	"bytes"
	"context"
	"encoding/json"
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
	s.createDefaultSeries()
}

func (s *APITestSuite) TearDownTest() {
	s.cleanAllTables()
}

func (s *APITestSuite) createDefaultSeries() {
	s.T().Helper()
	series := models.Series{
		ID:          0,
		Description: "default series desc",
		Title:       "default series title",
	}

	_, err := s.testDB.CreateSeries(context.Background(), series)
	s.Require().NoError(err)
}

func (s *APITestSuite) cleanAllTables() {
	s.T().Helper()

	for _, tableName := range []string{
		"forecaster.series",
		"forecaster.polls",
		"forecaster.options",
	} {
		_, err := s.testDB.DB.Exec(context.Background(), "TRUNCATE TABLE "+tableName+" CASCADE")
		s.Require().NoError(err)
	}
}

func (s *APITestSuite) CrudEndpointTest[T any](create T) {
	// create series
	cs := swagger.CreateSeries{
		Description: "test desc",
		Title:       "test title",
	}

	b, err := json.Marshal(cs)
	s.Require().NoError(err)

	s.Require().NoError(err)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/series", envs.AppPort),
		"application/json",
		bytes.NewReader(b))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
