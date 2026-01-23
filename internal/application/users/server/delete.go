package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-ddd-template/pkg/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	protobuf "go-ddd-template/generated/server"
	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h UserHandlers) DeleteUser(
	ctx context.Context,
	req *protobuf.DeleteUserRequest,
) (*emptypb.Empty, error) {
	userInfo := auth.GetUserInfo(ctx)
	if !userInfo.IsAdmin() {
		return nil, auth.ErrNoAdminRole
	}

	id, err := valueobjects.NewUserIDFromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	if err := h.userService.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to delete user",
			loggerutils.ErrAttr(fmt.Errorf("failed to delete user: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &emptypb.Empty{}, nil
}
