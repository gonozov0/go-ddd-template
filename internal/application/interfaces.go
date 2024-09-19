package application

import (
	"context"

	ordersDomain "go-echo-template/internal/domain/orders"
	productsDomain "go-echo-template/internal/domain/products"
	usersDomain "go-echo-template/internal/domain/users"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*usersDomain.User, error)
	SaveUser(ctx context.Context, u usersDomain.User) error
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, o *ordersDomain.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*ordersDomain.Order, error)
}

type ProductRepository interface {
	GetProductsForUpdate(ctx context.Context, ids []uuid.UUID) ([]productsDomain.Product, error)
	SaveProducts(ctx context.Context, ps []productsDomain.Product) error
}
