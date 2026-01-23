package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-ddd-template/pkg/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protobuf "go-ddd-template/generated/server"
	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h UserHandlers) GetUser(ctx context.Context, req *protobuf.GetUserRequest) (*protobuf.GetUserResponse, error) {
	userInfo := auth.GetUserInfo(ctx)
	if !userInfo.IsAdmin() {
		return nil, auth.ErrNoAdminRole
	}

	id, err := valueobjects.NewUserIDFromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}

		slog.ErrorContext(ctx, "failed to get user", loggerutils.ErrAttr(fmt.Errorf("failed to get user: %w", err)))

		return nil, grpcutils.ErrInternalError
	}

	return &protobuf.GetUserResponse{
		Id:    user.GetID().String(),
		Name:  user.GetName().String(),
		Email: user.GetEmail().String(),
	}, nil
}
