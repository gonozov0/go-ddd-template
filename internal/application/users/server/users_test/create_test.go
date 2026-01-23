package users_test

import (
	"google.golang.org/grpc/codes"

	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	"go-ddd-template/pkg/grpcutils"
)

func (s *UsersSuite) TestCreateUser() {
	s.Run("User created successfully", func() {
		user := userunithelpers.GenerateUser(s)

		resp, err := s.GRPCHandlers.CreateUser(s.AdminCtx, userunithelpers.ToCreateUserRequest(user))
		grpcutils.CheckCodeWithSuite(s, err, codes.OK)

		s.Require().Equal(user.GetID().String(), resp.Id)
	})
}
