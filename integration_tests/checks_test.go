package integrationtests

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/integration_tests/suites"
)

type APISuite struct {
	suites.BaseSuite
}

func (s *APISuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
}

func (s *APISuite) TestPing() {
	resp, err := http.Get(s.ServerURL + "/ping/")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NoError(resp.Body.Close())
}

func (s *APISuite) TestCheckErrors() {
	req, err := http.NewRequest(
		http.MethodGet,
		s.ServerURL+"/checks/errors",
		nil,
	)
	s.Require().NoError(err)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(resp.Body.Close())
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode, string(body))
}

func (s *APISuite) TestCheckPanics() {
	req, err := http.NewRequest(
		http.MethodGet,
		s.ServerURL+"/checks/panics",
		nil,
	)
	s.Require().NoError(err)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)

	s.Require().NoError(err)
	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(resp.Body.Close())
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode, string(body))
}

func (s *APISuite) TestSwagger() {
	resp, err := http.Get(s.ServerURL + "/swagger/")

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NoError(resp.Body.Close())
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
