package users_test

import (
	"testing"

	"go-echo-ddd-template/internal/application"
	"go-echo-ddd-template/internal/application/users"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	application.ServerSuite
	UserHandlers users.UserHandlers
}

func (s *UsersSuite) SetupSuite() {
	s.ServerSuite.SetupSuite()
	s.UserHandlers = users.SetupHandlers(s.UsersRepo)
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
