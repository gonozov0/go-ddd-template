package orders

import (
	"context"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/internal/service/users"
)

type OrderRepository interface {
	CreateOrder(
		ctx context.Context,
		productIDs valueobjects.ProductIDs,
		createFn func([]domain.Product) (*domain.Order, error),
	) (*domain.Order, error)
	GetOrder(ctx context.Context, id valueobjects.OrderID) (*domain.Order, error)
	DeleteOrder(ctx context.Context, orderID valueobjects.OrderID) error
	ProcessOrders(ctx context.Context, updateOrders func([]*domain.Order) error) error
}

type OrderService struct {
	orderRepo   OrderRepository
	userService users.UserService
}

func NewOrderService(or OrderRepository, us users.UserService) OrderService {
	return OrderService{
		orderRepo:   or,
		userService: us,
	}
}
