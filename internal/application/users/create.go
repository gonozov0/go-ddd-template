package users

import (
	"context"
	"errors"
	"net/http"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h UserHandlers) PostUsers(c echo.Context) error {
	var req openapi.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := users.CreateUser(req.Name, string(req.Email))
	if err != nil {
		msg := err.Error()
		if errors.Is(err, users.ErrInvalidUser) || errors.Is(err, users.ErrUserValidation) {
			return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Message: &msg})
		}
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: &msg})
	}

	if err := h.repo.SaveUser(*user); err != nil {
		msg := err.Error()
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: &msg})
	}

	id := user.ID()
	return c.JSON(http.StatusCreated, openapi.CreateUserResponse{Id: &id})
}

func (h UserHandlers) CreateUser(
	_ context.Context,
	req *protobuf.CreateUserRequest,
) (*protobuf.CreateUserResponse, error) {
	user, err := users.CreateUser(req.GetName(), req.GetEmail())
	if err != nil {
		if errors.Is(err, users.ErrInvalidUser) || errors.Is(err, users.ErrUserValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := h.repo.SaveUser(*user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protobuf.CreateUserResponse{
		Id: user.ID().String(),
	}, nil
}
