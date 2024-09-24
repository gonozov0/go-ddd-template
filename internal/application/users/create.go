package users

import (
	"context"
	"errors"
	"net/http"

	"go-echo-template/generated/openapi"
	"go-echo-template/generated/protobuf"
	domain "go-echo-template/internal/domain/users"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h UserHandlers) PostUsers(c echo.Context) error {
	ctx := c.Request().Context()
	var req openapi.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	email := string(req.Email)
	user, err := h.repo.CreateUser(ctx, email, func() (*domain.User, error) {
		return domain.CreateUser(req.Name, email)
	})
	if err != nil {
		msg := err.Error()
		if errors.Is(err, domain.ErrInvalidUser) || errors.Is(err, domain.ErrUserValidation) {
			return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Message: &msg})
		}
		if errors.Is(err, domain.ErrUserAlreadyExist) {
			return c.JSON(http.StatusConflict, openapi.ErrorResponse{Message: &msg})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: &msg})
	}

	id := user.ID()
	return c.JSON(http.StatusCreated, openapi.CreateUserResponse{Id: &id})
}

func (h UserHandlers) CreateUser(
	ctx context.Context,
	req *protobuf.CreateUserRequest,
) (*protobuf.CreateUserResponse, error) {
	email := req.GetEmail()
	user, err := h.repo.CreateUser(ctx, email, func() (*domain.User, error) {
		return domain.CreateUser(req.GetName(), email)
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) || errors.Is(err, domain.ErrUserValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protobuf.CreateUserResponse{
		Id: user.ID().String(),
	}, nil
}
