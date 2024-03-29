package users

import (
	"sync"

	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

type user struct {
	Name  string
	Email string
}

type InMemoryRepo struct {
	users map[uuid.UUID]user
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		users: make(map[uuid.UUID]user),
	}
}

func (r *InMemoryRepo) SaveUser(u users.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[u.ID()] = user{
		Name:  u.Name(),
		Email: u.Email(),
	}
	return nil
}

func (r *InMemoryRepo) GetUser(id uuid.UUID) (*users.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
