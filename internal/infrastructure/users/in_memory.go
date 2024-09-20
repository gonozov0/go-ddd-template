package users

import (
	"context"

	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type user struct {
	Name  string
	Email string
}

type InMemoryRepo struct {
	users map[uuid.UUID]user
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		users: make(map[uuid.UUID]user),
	}
}

func (r *InMemoryRepo) SaveUser(_ context.Context, u users.User) error {
	r.users[u.ID()] = user{
		Name:  u.Name(),
		Email: u.Email(),
	}
	return nil
}

func (r *InMemoryRepo) GetUser(_ context.Context, id uuid.UUID) (*users.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, users.ErrUserNotFound
	}
	user, err := users.NewUser(id, u.Name, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
