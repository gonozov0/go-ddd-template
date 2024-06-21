package users

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/generated/openapi"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
