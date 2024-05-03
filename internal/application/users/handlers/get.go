package handlers

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/internal/domain/users"
	"go-echo-ddd-template/pkg/responses"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUser retrieves a user by UUID from the system.
//
//	@Summary		Get a user by ID
//	@Description	Retrieves user details using the user ID provided in the path
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"User ID"
//	@Success		200	{object}	GetResponse				"User details successfully retrieved"
//	@Failure		400	{object}	responses.ErrorResponse	"Invalid user ID"
//	@Failure		404	{object}	responses.ErrorResponse	"User not found"
//	@Failure		500	{object}	responses.ErrorResponse	"Failed to get user"
//	@Router			/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse{Message: "Invalid user ID"})
	}

	user, err := h.repo.GetUser(uuidID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse{Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Message: "Failed to get user"})
	}

	return c.JSON(http.StatusOK, GetResponse{
		ID:    user.ID().String(),
		Name:  user.Name(),
		Email: user.Email(),
	})
}
