package handlers

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid user ID"})
	}

	user, err := h.repo.GetUser(uuidID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user"})
	}

	return c.JSON(http.StatusOK, GetResponse{
		ID:    user.ID().String(),
		Name:  user.Name(),
		Email: user.Email(),
	})
}
