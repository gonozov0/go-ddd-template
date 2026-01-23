package application_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/integration_tests/suites"
	userhelper "go-ddd-template/integration_tests/users/helpers"
)

type UsersSuite struct {
	suites.BaseSuite

	grpcClient pb.UserServiceClient

	userHelper *userhelper.HelperSuite
}

func (s *UsersSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.grpcClient = pb.NewUserServiceClient(s.ServerConn)

	s.userHelper = userhelper.NewHelperSuite(&s.BaseSuite)
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}
