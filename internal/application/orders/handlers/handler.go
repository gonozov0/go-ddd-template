package handlers

import (
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
}

type Handler struct {
	orderRepo   OrderRepository
	userRepo    UserRepository
	productRepo ProductRepository
}

func NewHandler(or OrderRepository, ur UserRepository, pr ProductRepository) *Handler {
	return &Handler{
		orderRepo:   or,
		userRepo:    ur,
		productRepo: pr,
	}
}
