package e2e_test

import (
	"net/http"

	"github.com/levigross/grequests"
)

func (suite *APITestSuite) TestPing() {
	resp, err := grequests.Get(suite.serverURL+"ping", nil)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
	suite.Require().Equal("pong", resp.String())
}
