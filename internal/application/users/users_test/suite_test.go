package users_test

import (
	"testing"

	"go-echo-template/internal/application"
	"go-echo-template/internal/application/users"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	application.ServerSuite
	GRPCHandlers users.UserHandlers
}

func (s *UsersSuite) SetupTest() {
	s.ServerSuite.SetupTest()
	s.GRPCHandlers = users.SetupHandlers(s.UsersRepo)
}

func TestUsersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UsersSuite))
}
