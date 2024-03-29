package orders

import (
	"errors"

	"go-echo-ddd-template/internal/domain/orders"

	"github.com/google/uuid"
)

type PostgresRepo struct {
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{}
}

func (r *PostgresRepo) SaveOrder(_ orders.Order) error {
	return errors.New("not implemented")
}

func (r *PostgresRepo) GetOrder(_ uuid.UUID) (*orders.Order, error) {
	return nil, errors.New("not implemented")
}
