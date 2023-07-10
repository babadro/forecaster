package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/suite"
)

type envVars struct {
	AppPort int `env:"APP_PORT,required"`
}

// APITestSuite defines the suite
type APITestSuite struct {
	suite.Suite

	client *http.Client
}

// SetupSuite function will be run by testify before any tests or test suites are run.
func (s *APITestSuite) SetupSuite() {
	var envs envVars
	s.Require().NoError(env.Parse(&envs))

	s.client = &http.Client{
		Timeout: time.Second * 10,
	}
}

func (s *APITestSuite) TearDownSuite() {

}

func (s *APITestSuite) SetupTest() {

}

func (s *APITestSuite) TearDownTest() {

}

func (s *APITestSuite) CreateDefaultSeries() {
	resp, err := s.client.Post()
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
