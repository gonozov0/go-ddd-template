package users

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/internal/domain/users"
)

func (s UserService) DeleteUser(ctx context.Context, id valueobjects.UserID) error {
	err := s.repo.DeleteUser(ctx, id, func(u *users.User) error {
		// must be implemented any checks and validations if needed
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
