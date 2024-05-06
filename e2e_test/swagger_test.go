package e2e_test

import (
	"net/http"

	"github.com/levigross/grequests"
)

func (suite *APITestSuite) TestSwagger() {
	resp, err := grequests.Get(suite.serverURL+"swagger", nil)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}
