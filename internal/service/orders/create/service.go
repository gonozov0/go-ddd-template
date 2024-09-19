package create

import (
	"context"

	"go-echo-template/internal/domain/orders"
	"go-echo-template/internal/domain/products"
	"go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type orderRepo interface {
	SaveOrder(ctx context.Context, o *orders.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}

type userRepo interface {
	GetUser(ctx context.Context, id uuid.UUID) (*users.User, error)
}

type productRepo interface {
	GetProductsForUpdate(ctx context.Context, ids []uuid.UUID) ([]products.Product, error)
	SaveProducts(ctx context.Context, ps []products.Product) error
}

type OrderCreationService struct {
	orderRepo   orderRepo
	userRepo    userRepo
	productRepo productRepo
}

func NewOrderCreationService(or orderRepo, ur userRepo, pr productRepo) *OrderCreationService {
	return &OrderCreationService{
		orderRepo:   or,
		userRepo:    ur,
		productRepo: pr,
	}
}
