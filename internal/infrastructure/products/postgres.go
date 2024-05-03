package products

import (
	"errors"

	"go-echo-ddd-template/internal/domain/products"

	"github.com/google/uuid"
)

type PostgresRepo struct {
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{}
}

func (r *PostgresRepo) SaveProducts(_ []products.Product) error {
	return errors.New("not implemented")
}

func (r *PostgresRepo) GetProductsForUpdate(_ []uuid.UUID) ([]products.Product, error) {
	return nil, errors.New("not implemented")
}

func (r *PostgresRepo) CancelUpdate() {}
