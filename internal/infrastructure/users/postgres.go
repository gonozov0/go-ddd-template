package users

import (
	"errors"

	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type PostgresRepo struct {
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{}
}

func (r *PostgresRepo) SaveUser(_ *users.User) error {
	return errors.New("not implemented")
}

func (r *PostgresRepo) GetUser(_ uuid.UUID) (*users.User, error) {
	return nil, errors.New("not implemented")
}
