package application_test

import (
	"testing"

	"go-echo-template/generated/openapi"
	"go-echo-template/generated/protobuf"
	"go-echo-template/integration_tests/apisuite"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
	apisuite.APITestSuite

	grpcClient protobuf.UserServiceClient
	httpClient *openapi.Client
}

func (s *UsersSuite) SetupSuite() {
	s.APITestSuite.SetupSuite("8082")
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
	t.Parallel()
	suite.Run(t, new(UsersSuite))
}
