package users

import (
	"context"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
)

type UserRepository interface {
	CreateUser(ctx context.Context, createFn func() (*domain.User, error)) (*domain.User, error)
	UpdateUser(ctx context.Context, id valueobjects.UserID, updateFn func(*domain.User) error) (*domain.User, error)
	GetUser(ctx context.Context, id valueobjects.UserID) (*domain.User, error)
	DeleteUser(ctx context.Context, id valueobjects.UserID, deleteFn func(*domain.User) error) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) UserService {
	return UserService{
		repo: r,
	}
}
