package users_test

import (
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"

	pb "go-ddd-template/generated/server"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/grpcutils"
)

func (s *UsersSuite) TestGetUser() {
	user := userunithelpers.CreateUser(s, s.UsersRepo)

	resp, err := s.GRPCHandlers.GetUser(s.AdminCtx, &pb.GetUserRequest{Id: user.GetID().String()})
	grpcutils.CheckCodeWithSuite(s, err, codes.OK)

	s.Require().Equal(user.GetID().String(), resp.GetId())
	s.Require().Equal(user.GetName().String(), resp.GetName())
	s.Require().Equal(user.GetEmail().String(), resp.GetEmail())
}

func (s *UsersSuite) TestGetUserNotFound() {
	_, err := s.GRPCHandlers.GetUser(s.AdminCtx, &pb.GetUserRequest{Id: uuid.NewString()})
	message := grpcutils.CheckCodeWithSuite(s, err, codes.NotFound)

	s.Require().Contains(message, domain.ErrUserNotFound.Error())
}
