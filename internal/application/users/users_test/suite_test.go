package users_test

import (
	"testing"

	"go-echo-ddd-template/internal/application"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	application.ServerSuite
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
