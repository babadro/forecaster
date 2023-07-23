package polls_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/babadro/forecaster/internal/models/swagger"
	"net/http"
)

func (s *APITestSuite) TestSeries() {
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
