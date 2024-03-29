package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go-echo-ddd-template/internal/application/users/handlers"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

func (s *UsersSuite) TestGetUser() {
	user, _ := users.CreateUser("test", "test@test.com")
	err := s.repo.SaveUser(*user)
	s.Require().NoError(err)

	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID().String(), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var resp handlers.GetResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().Equal(handlers.GetResponse{
		ID:    user.ID().String(),
		Name:  user.Name(),
		Email: user.Email(),
	}, resp)
}

func (s *UsersSuite) TestGetUserNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusNotFound, rec.Code, rec.Body.String())
	s.Require().Equal(`{"message":"User not found"}`+"\n", rec.Body.String())
}
