package users

import (
	"context"

	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
)

type user struct {
	Name  string
	Email string
}

type InMemoryRepo struct {
	users map[valueobjects.UserID]user
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		users: make(map[valueobjects.UserID]user),
	}
}

func (r *InMemoryRepo) CreateUser(
	_ context.Context,
	createFn func() (*domain.User, error),
) (*domain.User, error) {
	u, err := createFn()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if r.checkUserExist(u.GetEmail().String()) {
		return nil, domain.ErrUserAlreadyExist
	}

	r.users[u.GetID()] = user{
		Name:  u.GetName().String(),
		Email: u.GetEmail().String(),
	}

	return u, nil
}

func (r *InMemoryRepo) UpdateUser(
	_ context.Context,
	id valueobjects.UserID,
	updateFn func(*domain.User) error,
) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	userName, err := domain.NewName(u.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init user name: %w", err)
	}

	userEmail, err := valueobjects.NewEmail(u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to init user email: %w", err)
	}

	entity := domain.NewUser(id, userName, userEmail)

	err = updateFn(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	r.users[id] = user{
		Name:  entity.GetName().String(),
		Email: entity.GetEmail().String(),
	}

	return entity, nil
}

func (r *InMemoryRepo) GetUser(_ context.Context, id valueobjects.UserID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	userName, err := domain.NewName(u.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init user name: %w", err)
	}

	userEmail, err := valueobjects.NewEmail(u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to init user email: %w", err)
	}

	user := domain.NewUser(id, userName, userEmail)

	return user, nil
}

func (r *InMemoryRepo) DeleteUser(_ context.Context, id valueobjects.UserID, deleteFn func(*domain.User) error) error {
	u, ok := r.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}

	userName, err := domain.NewName(u.Name)
	if err != nil {
		return fmt.Errorf("failed to init user name: %w", err)
	}

	userEmail, err := valueobjects.NewEmail(u.Email)
	if err != nil {
		return fmt.Errorf("failed to init user email: %w", err)
	}

	user := domain.NewUser(id, userName, userEmail)

	if err := deleteFn(user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	delete(r.users, id)

	return nil
}

func (r *InMemoryRepo) checkUserExist(email string) bool {
	for _, u := range r.users {
		if u.Email == email {
			return true
		}
	}

	return false
}
