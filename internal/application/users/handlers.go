package users

import (
	"context"

	"go-echo-template/generated/protobuf"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, email string, createFn func() (*users.User, error)) (*users.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(*users.User) (bool, error)) (*users.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}

type UserHandlers struct {
	protobuf.UnimplementedUserServiceServer
	repo Repository
}

func SetupHandlers(repo Repository) UserHandlers {
	return UserHandlers{
		repo: repo,
	}
}
