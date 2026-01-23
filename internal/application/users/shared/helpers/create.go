package helpers

import (
	"context"

	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/testify"
)

type UserCreater interface {
	CreateUser(
		ctx context.Context,
		createFn func() (*domain.User, error),
	) (*domain.User, error)
}

func CreateUser(s testify.Suite, repo UserCreater, opts ...GenerateUserOption) *domain.User {
	user := GenerateUser(s, opts...)

	_, err := repo.CreateUser(context.Background(), func() (*domain.User, error) { return user, nil })
	s.Require().NoError(err)

	return user
}
