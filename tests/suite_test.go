package polls_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/babadro/forecaster/tests/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

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

type crudEndpointTestInput[C, R, U any] struct {
	createInput    C
	updateInput    U
	checkCreateRes func(t *testing.T, got R)
	checkReadRes   func(t *testing.T, got R)
	checkUpdateRes func(t *testing.T, expectedID int32, got R)
	path           string
}

func testCRUDEndpoints[C, R, U any](t *testing.T, in crudEndpointTestInput[C, R, U]) {
	// create
	gotCreateResult := create[C, R](t, in.createInput, in.path)
	in.checkCreateRes(t, gotCreateResult)

	itemID := id(gotCreateResult)

	// read
	readResp, err := http.Get(
		fmt.Sprintf("http://localhost:%d/%s/%d", envs.AppPort, in.path, itemID),
	)
	require.NoError(t, err)
	defer func() { _ = readResp.Body.Close() }()

	require.Equal(t, http.StatusOK, readResp.StatusCode)
	in.checkReadRes(t, gotCreateResult)

	// update
	b, err := json.Marshal(in.updateInput)
	require.NoError(t, err)

	updateReq, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://localhost:%d/%s/%d", envs.AppPort, in.path, itemID),
		bytes.NewReader(b))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := http.DefaultClient.Do(updateReq)
	require.NoError(t, err)
	defer func() { _ = updateResp.Body.Close() }()

	require.Equal(t, http.StatusOK, updateResp.StatusCode)

	var gotUpdateResult R
	err = json.NewDecoder(updateResp.Body).Decode(&gotUpdateResult)
	require.NoError(t, err)

	in.checkUpdateRes(t, itemID, gotUpdateResult)

	// delete
	deleteReq, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://localhost:%d/%s/%d", envs.AppPort, in.path, itemID),
		nil)
	require.NoError(t, err)

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	require.NoError(t, err)
	defer func() { _ = deleteResp.Body.Close() }()

	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	// read deleted
	readResp, err = http.Get(
		fmt.Sprintf("http://localhost:%d/%s/%d", envs.AppPort, in.path, itemID),
	)
	require.NoError(t, err)
	defer func() { _ = readResp.Body.Close() }()
	require.Equal(t, http.StatusNotFound, readResp.StatusCode)
}

func timeRoundEqualNow(t *testing.T, got strfmt.DateTime) {
	now := time.Now()
	require.True(t, now.Sub(time.Time(got)).Abs() < time.Second, "expected %v, got %v", now, got)
}

func timeRoundEqual(t *testing.T, expected, got strfmt.DateTime) {
	require.True(t, time.Time(expected).Sub(time.Time(got)).Abs() < time.Second,
		"expected %v, got %v", expected, got)
}

func id(entity interface{}) int32 {
	switch v := entity.(type) {
	case models.Series:
		return v.ID
	case models.Poll:
		return v.ID
	case models.Option:
		return v.ID
	default:
		panic(fmt.Sprintf("unknown type %T", entity))
	}
}

func create[IN any, OUT any](t *testing.T, in IN, path string) OUT {
	b, err := json.Marshal(in)
	require.NoError(t, err)

	createResp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/%s", envs.AppPort, path),
		"application/json",
		bytes.NewReader(b))
	require.NoError(t, err)

	defer func() { _ = createResp.Body.Close() }()

	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var gotCreateResult OUT
	err = json.NewDecoder(createResp.Body).Decode(&gotCreateResult)
	require.NoError(t, err)

	return gotCreateResult
}

func randomModel[T any](t *testing.T) T {
	var model T
	require.NoError(t, gofakeit.Struct(&model))

	return model
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
