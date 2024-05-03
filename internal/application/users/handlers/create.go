package handlers

import (
	"errors"
	"net/http"

	"go-echo-ddd-template/internal/domain/users"
	"go-echo-ddd-template/pkg/responses"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateResponse struct {
	// ID is the UUID of the newly created user.
	ID uuid.UUID `json:"id"`
}

// CreateUser creates a new user in the system.
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided name and email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateRequest			true	"User creation request"
//	@Success		201		{object}	CreateResponse			"User successfully created"
//	@Failure		400		{object}	responses.ErrorResponse	"Invalid input data"
//	@Failure		500		{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var data CreateRequest
	if err := c.Bind(&data); err != nil {
		return err
	}

	user, err := users.CreateUser(data.Name, data.Email)
	if err != nil {
		if errors.Is(err, users.ErrInvalidUser) || errors.Is(err, users.ErrUserValidation) {
			return c.JSON(http.StatusBadRequest, responses.ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Message: err.Error()})
	}

	if err := h.repo.SaveUser(*user); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, CreateResponse{ID: user.ID()})
}
