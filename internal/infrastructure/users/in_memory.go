package users

import (
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

func (r *InMemoryRepo) SaveUser(u *users.User) error {
	if u == nil {
		return users.ErrInvalidUser
	}
	r.users[u.GetID()] = user{
		Name:  u.GetName(),
		Email: u.GetEmail(),
	}
	return nil
}

func (r *InMemoryRepo) GetUser(id uuid.UUID) (*users.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, users.ErrUserNotFound
	}
	user, err := users.NewUserWithID(id, u.Name, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
