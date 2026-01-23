package users

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
)

type UserToCreate struct {
	ID    valueobjects.UserID
	Name  domain.Name
	Email valueobjects.Email
}

func (s UserService) CreateUser(ctx context.Context, userToCreate UserToCreate) (*domain.User, error) {
	user, err := s.repo.CreateUser(ctx, func() (*domain.User, error) {
		return domain.NewUser(userToCreate.ID, userToCreate.Name, userToCreate.Email), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
