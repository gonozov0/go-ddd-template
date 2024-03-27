package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go-echo-template/internal/application/users/handlers"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

func (s *UsersSuite) TestGetUser() {
	user, _ := users.NewUser("test", "test@test.com")
	err := s.repo.SaveUser(user)
	s.Require().NoError(err)

	req := httptest.NewRequest(http.MethodGet, "/users/"+user.GetID().String(), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var resp handlers.GetResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().Equal(handlers.GetResponse{
		ID:    user.GetID().String(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}, resp)
}

func (s *UsersSuite) TestGetUserNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusNotFound, rec.Code, rec.Body.String())
	s.Require().Equal(`{"message":"User not found"}`+"\n", rec.Body.String())
}
