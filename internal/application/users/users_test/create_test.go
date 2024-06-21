package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go-echo-ddd-template/generated/openapi"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *UsersSuite) TestCreateUserSuccess() {
	userReq := openapi.CreateUserRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	reqBody, _ := json.Marshal(userReq)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.Echo.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusCreated, rec.Code)
	var resp openapi.CreateUserResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotEqual(resp.Id, uuid.Nil)
}
