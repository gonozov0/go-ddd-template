package orders

import (
	"go-echo-ddd-template/generated/protobuf"
	"go-echo-ddd-template/internal/domain/orders"
	"go-echo-ddd-template/internal/domain/products"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

type OrderRepository interface {
	SaveOrder(o orders.Order) error
	GetOrder(id uuid.UUID) (*orders.Order, error)
}

type UserRepository interface {
	GetUser(id uuid.UUID) (*users.User, error)
}

type ProductRepository interface {
	GetProductsForUpdate(ids []uuid.UUID) ([]products.Product, error)
	SaveProducts(ps []products.Product) error
	CancelUpdate()
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
