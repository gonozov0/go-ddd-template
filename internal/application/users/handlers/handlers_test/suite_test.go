package handlers_test

import (
	"testing"

	"go-echo-ddd-template/internal/application/users"
	"go-echo-ddd-template/internal/application/users/handlers"
	infra "go-echo-ddd-template/internal/infrastructure/users"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	e    *echo.Echo
	repo handlers.Repository
}

// SetupSuite is used to initialize resources before running the all tests in suite.
func (s *UsersSuite) SetupSuite() {
	s.e = echo.New()
	s.repo = infra.NewInMemoryRepo()
	users.Setup(s.e, s.repo)
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
