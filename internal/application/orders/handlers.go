package orders

import (
	"context"

	"go-echo-template/generated/protobuf"
	"go-echo-template/internal/domain/orders"
	"go-echo-template/internal/domain/products"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, o *orders.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}

type ProductRepository interface {
	GetProductsForUpdate(ctx context.Context, ids []uuid.UUID) ([]products.Product, error)
	SaveProducts(ctx context.Context, ps []products.Product) error
}

type OrderHandlers struct {
	protobuf.UnimplementedOrderServiceServer
	orderRepo   OrderRepository
	userRepo    UserRepository
	productRepo ProductRepository
}

func SetupHandlers(or OrderRepository, ur UserRepository, pr ProductRepository) OrderHandlers {
	return OrderHandlers{
		orderRepo:   or,
		userRepo:    ur,
		productRepo: pr,
	}
}
