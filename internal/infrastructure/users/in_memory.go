package users

import (
	"context"
	"fmt"

	domain "go-echo-template/internal/domain/users"

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

func (r *InMemoryRepo) CreateUser(
	_ context.Context,
	email string,
	createFn func() (*domain.User, error),
) (*domain.User, error) {
	if r.checkUserExist(email) {
		return nil, domain.ErrUserAlreadyExist
	}

	u, err := createFn()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.users[u.ID()] = user{
		Name:  u.Name(),
		Email: email,
	}
	return u, nil
}

func (r *InMemoryRepo) UpdateUser(
	_ context.Context,
	id uuid.UUID,
	updateFn func(*domain.User) (bool, error),
) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	entity, err := domain.NewUser(id, u.Name, u.Email)
	if err != nil {
		return nil, err
	}
	updated, err := updateFn(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if !updated {
		return entity, nil
	}

	r.users[id] = user{
		Name:  entity.Name(),
		Email: entity.Email(),
	}
	return entity, nil
}

func (r *InMemoryRepo) SaveUser(_ context.Context, u domain.User) error {
	r.users[u.ID()] = user{
		Name:  u.Name(),
		Email: u.Email(),
	}
	return nil
}

func (r *InMemoryRepo) GetUser(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	user, err := domain.NewUser(id, u.Name, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *InMemoryRepo) checkUserExist(email string) bool {
	for _, u := range r.users {
		if u.Email == email {
			return true
		}
	}
	return false
}
