package users

import (
	"context"
	"errors"
	"net/http"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h UserHandlers) GetUsersId( //nolint:revive,stylecheck // fit to generated code
	c echo.Context,
	id openapi_types.UUID,
) error {
	user, err := h.repo.GetUser(id)
	if err != nil {
		msg := err.Error()
		if errors.Is(err, users.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, openapi.ErrorResponse{Message: &msg})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: &msg})
	}

	name := user.Name()
	email := openapi_types.Email(user.Email())
	return c.JSON(http.StatusOK, openapi.GetUserResponse{
		Id:    &id,
		Name:  &name,
		Email: &email,
	})
}

func (h UserHandlers) GetUser(_ context.Context, req *protobuf.GetUserRequest) (*protobuf.GetUserResponse, error) {
	uid, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID")
	}
	user, err := h.repo.GetUser(uid)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &protobuf.GetUserResponse{
		Id:    user.ID().String(),
		Name:  user.Name(),
		Email: user.Email(),
	}, nil
}
