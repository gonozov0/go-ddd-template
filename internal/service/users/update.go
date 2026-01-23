package users

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
)

type UserToUpdate struct {
	ID    valueobjects.UserID
	Name  domain.Name
	Email valueobjects.Email
}

func (s UserService) UpdateUser(
	ctx context.Context,
	userToUpdate UserToUpdate,
) (*domain.User, error) {
	user, err := s.repo.UpdateUser(ctx, userToUpdate.ID, func(user *domain.User) error {
		user.Update(userToUpdate.Name, userToUpdate.Email)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
