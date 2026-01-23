package users

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/internal/domain/users"
)

func (s UserService) GetUser(ctx context.Context, id valueobjects.UserID) (*users.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
