package polls_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

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
	AppPort   int    `env:"APP_PORT" envDefault:"8080"`
	DBConn    string `env:"DB_CONN" envDefault:"postgres://postgres:postgres@localhost:5432/forecaster?sslmode=disable"`
	SleepMode bool   `env:"SLEEP_MODE" envDefault:"false"`
}

// APITestSuite defines the suite...
type APITestSuite struct {
	suite.Suite

	testDB *db.TestDB

	apiAddr string
	client  *http.Client

	sleepMode bool
}

// SetupSuite function will be run by testify before any tests or test suites are run.
func (s *APITestSuite) SetupSuite() {
	var envs envVars

	s.Require().NoError(env.Parse(&envs))

	s.apiAddr = fmt.Sprintf("http://localhost:%d", envs.AppPort)

	s.sleepMode = envs.SleepMode

	s.client = &http.Client{
		Timeout: time.Second * 10,
	}

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}

	s.testDB = db.NewTestDB(dbPool)
}

func (s *APITestSuite) TearDownSuite() {

}

func (s *APITestSuite) SetupTest() {
	s.createDefaultSeries()
}

func (s *APITestSuite) TearDownTest() {
	if s.sleepMode {
		time.Sleep(time.Hour * 10_000)
	}

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

	s.Require().NoError(s.testDB.CleanAllTables(context.Background()))
}

func (s *APITestSuite) url(path string) string {
	return fmt.Sprintf("%s/%s", s.apiAddr, path)
}

type crudEndpointTestInput[CIn, COut, R, UIn, UOut any] struct {
	createInput    CIn
	updateInput    UIn
	checkCreateRes func(t *testing.T, got COut)
	checkReadRes   func(t *testing.T, got R)
	checkUpdateRes func(t *testing.T, expectedID int32, got UOut)
	path           string
}

func testCRUDEndpoints[CIn, COut, R, UIn, UOut any](
	t *testing.T, in crudEndpointTestInput[CIn, COut, R, UIn, UOut], apiAddr string,
) {
	// create
	gotCreateResult := create[CIn, COut](t, in.createInput, apiAddr+"/"+in.path)
	in.checkCreateRes(t, gotCreateResult)

	itemID := id(gotCreateResult)

	// read
	gotReadResult := read[R](t, urlWithID(apiAddr, in.path, itemID))
	in.checkReadRes(t, gotReadResult)

	// update
	gotUpdateResult := update[UIn, UOut](t, in.updateInput, urlWithID(apiAddr, in.path, itemID))
	in.checkUpdateRes(t, itemID, gotUpdateResult)

	// delete
	deleteOp(t, urlWithID(apiAddr, in.path, itemID))

	// read deleted
	readShouldNotFound(t, urlWithID(apiAddr, in.path, itemID))
}

func timeRoundEqualNow(t *testing.T, got strfmt.DateTime) {
	now := time.Now()
	require.True(t, now.Sub(time.Time(got)).Abs() < time.Second, "now: %v, got: %v", now, got)
}

func timeRoundEqual(t *testing.T, expected, got strfmt.DateTime) {
	require.True(t, time.Time(expected).Sub(time.Time(got)).Abs() < time.Second,
		"now: %v, got: %v", expected, got)
}

func id(entity interface{}) int32 {
	switch v := entity.(type) {
	case models.Series:
		return v.ID
	case models.Poll:
		return v.ID
	default:
		panic(fmt.Sprintf("unknown type %T", entity))
	}
}

func create[IN any, OUT any](t *testing.T, in IN, url string) OUT {
	b, err := json.Marshal(in)
	require.NoError(t, err)

	createResp, err := http.Post(
		url,
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

func read[OUT any](t *testing.T, url string) OUT {
	readResp, err := http.Get(url)
	require.NoError(t, err)

	defer func() { _ = readResp.Body.Close() }()

	require.Equal(t, http.StatusOK, readResp.StatusCode)

	var got OUT
	err = json.NewDecoder(readResp.Body).Decode(&got)
	require.NoError(t, err)

	return got
}

func readShouldNotFound(t *testing.T, url string) {
	readResp, err := http.Get(url)
	require.NoError(t, err)

	defer func() { _ = readResp.Body.Close() }()

	require.Equal(t, http.StatusNotFound, readResp.StatusCode)
}

func update[IN any, OUT any](t *testing.T, in IN, url string) OUT {
	b, err := json.Marshal(in)
	require.NoError(t, err)

	updateReq, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := http.DefaultClient.Do(updateReq)
	require.NoError(t, err)

	defer func() { _ = updateResp.Body.Close() }()

	require.Equal(t, http.StatusOK, updateResp.StatusCode)

	var got OUT
	err = json.NewDecoder(updateResp.Body).Decode(&got)
	require.NoError(t, err)

	return got
}

func updateShouldReturnError[IN any](t *testing.T, in IN, url string, expectedStatus int) models.Error {
	b, err := json.Marshal(in)
	require.NoError(t, err)

	updateReq, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := http.DefaultClient.Do(updateReq)
	require.NoError(t, err)

	defer func() { _ = updateResp.Body.Close() }()

	require.Equal(t, expectedStatus, updateResp.StatusCode)

	var got models.Error
	err = json.NewDecoder(updateResp.Body).Decode(&got)
	require.NoError(t, err)

	return got
}

func deleteOp(t *testing.T, url string) {
	deleteReq, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	require.NoError(t, err)

	defer func() { _ = deleteResp.Body.Close() }()

	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
}

func randomModel[T any](t *testing.T) T {
	var model T

	require.NoError(t, gofakeit.Struct(&model))

	return model
}

func urlWithID(apiAddr string, path string, id int32) string {
	return fmt.Sprintf("%s/%s/%d", apiAddr, path, id)
}

func optionURLWithIDs(apiAddr string, pollID int32, optionID int16) string {
	return fmt.Sprintf("%s/options/%d/%d", apiAddr, pollID, optionID)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
