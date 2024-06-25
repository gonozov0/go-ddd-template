package application_test

import (
	"testing"

	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/integration_tests/apisuite"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	apisuite.APITestSuite

	client protobuf.UserServiceClient
}

func (s *UsersSuite) SetupSuite() {
	s.APITestSuite.SetupSuite()
	s.client = protobuf.NewUserServiceClient(s.Conn)
}

func (s *UsersSuite) TearDownSuite() {
	s.APITestSuite.TearDownSuite()
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
