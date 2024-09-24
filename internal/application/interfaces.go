package application

import (
	"context"

	"go-echo-template/internal/domain/orders"
	"go-echo-template/internal/domain/products"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email string, createFn func() (*users.User, error)) (*users.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(*users.User) (bool, error)) (*users.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, itemIDs []uuid.UUID, createFn func() (*orders.Order, error)) (*orders.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}

type ProductRepository interface {
	GetProducts(ctx context.Context, ids []uuid.UUID) ([]products.Product, error)
}
