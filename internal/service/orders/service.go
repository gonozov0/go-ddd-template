package orders

import (
	"context"

	"go-echo-template/internal/domain/orders"
	"go-echo-template/internal/domain/products"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, itemIDs []uuid.UUID, createFn func() (*orders.Order, error)) (*orders.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}

type ProductRepository interface {
	GetProducts(ctx context.Context, ids []uuid.UUID) ([]products.Product, error)
}

type OrderCreationService struct {
	orderRepo   OrderRepository
	userRepo    UserRepository
	productRepo ProductRepository
}

func NewOrderCreationService(or OrderRepository, ur UserRepository, pr ProductRepository) *OrderCreationService {
	return &OrderCreationService{
		orderRepo:   or,
		userRepo:    ur,
		productRepo: pr,
	}
}
