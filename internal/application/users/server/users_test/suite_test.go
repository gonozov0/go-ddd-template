package users_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go-ddd-template/internal/application"
	users "go-ddd-template/internal/application/users/server"
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
