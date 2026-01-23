package helpers

import (
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/integration_tests/suites"
	userunithelpers "go-ddd-template/internal/application/users/shared/helpers"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type HelperSuite struct {
	baseSuite *suites.BaseSuite

	grpcClient pb.UserServiceClient
}

func NewHelperSuite(baseSuite *suites.BaseSuite) *HelperSuite {
	handlerSuite := &HelperSuite{
		baseSuite: baseSuite,

		grpcClient: pb.NewUserServiceClient(baseSuite.ServerConn),
	}

	return handlerSuite
}

func (s *HelperSuite) CreateUser(
	opts ...userunithelpers.GenerateUserOption,
) valueobjects.UserID {
	user := userunithelpers.GenerateUser(s.baseSuite, opts...)

	_, err := s.grpcClient.CreateUser(
		s.baseSuite.AdminCtx,
		userunithelpers.ToCreateUserRequest(user),
	)
	grpcutils.CheckCodeWithSuite(s.baseSuite, err, codes.OK)

	s.baseSuite.T().Cleanup(func() { s.DeleteUserInCleanup(user.GetID()) })

	return user.GetID()
}

func (s *HelperSuite) DeleteUserInCleanup(
	userID valueobjects.UserID,
) {
	if userID.IsEmpty() {
		return
	}

	_, err := s.grpcClient.DeleteUser(s.baseSuite.AdminCtx, &pb.DeleteUserRequest{
		Id: userID.String(),
	})

	grpcStatus, ok := status.FromError(err)
	if err != nil || !ok || grpcStatus.Code() != codes.OK {
		slog.Error(
			"failed to delete user",
			loggerutils.ErrAttr(err),
			loggerutils.Attr("grpc_status", grpcStatus.Code().String()),
			loggerutils.Attr("grpc_message", grpcStatus.Message()),
			loggerutils.Attr("user_id", userID.String()),
		)

		return
	}

	slog.Info("Deleted user", loggerutils.Attr("user_id", userID.String()))
}
