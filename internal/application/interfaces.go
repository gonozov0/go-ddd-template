package application

import (
	ordersDomain "go-echo-ddd-template/internal/domain/orders"
	productsDomain "go-echo-ddd-template/internal/domain/products"
	usersDomain "go-echo-ddd-template/internal/domain/users"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(id uuid.UUID) (*usersDomain.User, error)
	SaveUser(u usersDomain.User) error
}

type OrderRepository interface {
	SaveOrder(o ordersDomain.Order) error
	GetOrder(id uuid.UUID) (*ordersDomain.Order, error)
}

type ProductRepository interface {
	GetProductsForUpdate(ids []uuid.UUID) ([]productsDomain.Product, error)
	SaveProducts(ps []productsDomain.Product) error
	CancelUpdate()
}
