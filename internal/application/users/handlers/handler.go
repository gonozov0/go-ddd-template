package handlers

import (
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type Repository interface {
	SaveUser(u *users.User) error
	GetUser(id uuid.UUID) (*users.User, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
