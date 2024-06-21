package users

import (
	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

type Repository interface {
	SaveUser(u users.User) error
	GetUser(id uuid.UUID) (*users.User, error)
}

type UserHandlers struct {
	repo Repository
}

func SetupHandlers(repo Repository) UserHandlers {
	return UserHandlers{
		repo: repo,
	}
}
