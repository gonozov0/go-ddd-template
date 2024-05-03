package create

import (
	"go-echo-ddd-template/internal/domain/orders"
	"go-echo-ddd-template/internal/domain/products"
	"go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

type orderRepo interface {
	SaveOrder(o orders.Order) error
	GetOrder(id uuid.UUID) (*orders.Order, error)
}

type userRepo interface {
	GetUser(id uuid.UUID) (*users.User, error)
}

type productRepo interface {
	GetProductsForUpdate(ids []uuid.UUID) ([]products.Product, error)
	SaveProducts(ps []products.Product) error
	CancelUpdate()
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
