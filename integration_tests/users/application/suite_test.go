package application_test

import (
	"testing"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/integration_tests/apisuite"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	apisuite.APITestSuite

	grpcClient protobuf.UserServiceClient
	httpClient *openapi.Client
}

func (s *UsersSuite) SetupSuite() {
	s.APITestSuite.SetupSuite()
	s.grpcClient = protobuf.NewUserServiceClient(s.Conn)
	var err error
	s.httpClient, err = openapi.NewClient(s.HTTPServerURL)
	if err != nil {
		s.Fail("Failed to create HTTP client", err)
	}
}

func (s *UsersSuite) TearDownSuite() {
	s.APITestSuite.TearDownSuite()
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
