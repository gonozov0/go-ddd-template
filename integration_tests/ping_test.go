package integration_tests_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-echo-template/integration_tests/suites"

	"github.com/levigross/grequests"
)

type TestSuite struct {
	suite.Suite
	suites.RunServerSuite
}

func (suite *TestSuite) SetupSuite() {
	suite.RunServerSuite.SetupSuite("8081")
}

func (suite *TestSuite) TestPing() {
	resp, err := grequests.Get(suite.HTTPServerURL+"/ping", nil)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
	suite.Require().Equal("pong", resp.String())
}

func TestAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TestSuite))
}
