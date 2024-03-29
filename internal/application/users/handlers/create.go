package handlers

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/internal/domain/users"

	"github.com/labstack/echo/v4"
)

type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

// CreateUser brings some json data from request body and returns it as response.
func (h *Handler) CreateUser(c echo.Context) error {
	var data CreateRequest
	if err := c.Bind(&data); err != nil {
		return err
	}

	user, err := users.CreateUser(data.Name, data.Email)
	if err != nil {
		if errors.Is(err, users.ErrInvalidUser) || errors.Is(err, users.ErrUserValidation) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	if err := h.repo.SaveUser(*user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, CreateResponse{ID: user.ID().String()})
}
