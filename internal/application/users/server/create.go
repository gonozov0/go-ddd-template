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
	service "go-ddd-template/internal/service/users"
	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h UserHandlers) CreateUser(
	ctx context.Context,
	req *protobuf.CreateUserRequest,
) (*protobuf.CreateUserResponse, error) {
	userInfo := auth.GetUserInfo(ctx)
	if !userInfo.IsAdmin() {
		return nil, auth.ErrNoAdminRole
	}

	id, err := valueobjects.NewUserIDFromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	name, err := domain.NewName(req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	email, err := valueobjects.NewEmail(req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	user, err := h.userService.CreateUser(ctx, service.UserToCreate{
		ID:    id,
		Name:  name,
		Email: email,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) || errors.Is(err, domain.ErrUserValidation) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		slog.ErrorContext(
			ctx,
			"failed to create user",
			loggerutils.ErrAttr(fmt.Errorf("failed to create user: %w", err)),
		)

		return nil, grpcutils.ErrInternalError
	}

	return &protobuf.CreateUserResponse{
		Id: user.GetID().String(),
	}, nil
}
