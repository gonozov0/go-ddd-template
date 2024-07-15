package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/generated/protobuf"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *UsersSuite) TestCreateUser() {
	s.Run("HTTP", func() {
		userReq := openapi.CreateUserRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		}
		reqBody, _ := json.Marshal(userReq)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.HTTPServer.ServeHTTP(rec, req)

		s.Require().Equal(http.StatusCreated, rec.Code)
		var resp openapi.CreateUserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Require().NoError(err)
		s.Require().NotEqual("", resp.Id)
		s.Require().NotEqual(uuid.Nil, resp.Id)
	})

	s.Run("GRPC", func() {
		userReq := protobuf.CreateUserRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		}
		resp, err := s.UserHandlers.CreateUser(context.Background(), &userReq)

		s.Require().NoError(err)
		s.Require().NotEqual("", resp.GetId())
		s.Require().NotEqual(uuid.Nil, resp.GetId())
	})
}
