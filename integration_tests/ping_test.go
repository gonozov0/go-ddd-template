package integration_tests_test

import (
	"net/http"
	"testing"

	"go-echo-template/integration_tests/apisuite"

	"github.com/levigross/grequests"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	apisuite.APITestSuite
}

func (suite *TestSuite) SetupSuite() {
	suite.APITestSuite.SetupSuite("8081")
}

func (suite *TestSuite) TestPing() {
	resp, err := grequests.Get(suite.HTTPServerURL+"/ping", nil)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
	suite.Require().Equal("pong", resp.String())
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
