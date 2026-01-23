package helpers

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/testify"
)

// Методы для генерации Product
type (
	productToCreate struct {
		id     valueobjects.ProductID
		name   valueobjects.ProductName
		price  valueobjects.ProductPrice
		status valueobjects.ProductStatus
	}

	GenerateProductOption func(*productToCreate) error
)

func ProductWithID(id valueobjects.ProductID) GenerateProductOption {
	return func(product *productToCreate) error {
		product.id = id
		return nil
	}
}

func ProductWithName(name valueobjects.ProductName) GenerateProductOption {
	return func(product *productToCreate) error {
		product.name = name
		return nil
	}
}

func ProductWithPrice(price valueobjects.ProductPrice) GenerateProductOption {
	return func(product *productToCreate) error {
		product.price = price
		return nil
	}
}

func ProductWithStatus(status valueobjects.ProductStatus) GenerateProductOption {
	return func(product *productToCreate) error {
		product.status = status
		return nil
	}
}

// GenerateProduct генерирует Product с помощью переданных опций
// Если опции не переданы, генерируется Product со случайными данными:
//   - id	  - случайный UUID
//   - name	  - случайное название продукта
//   - price  - случайная цена от 1.0 до 1000.0
//   - status - статус Init
func GenerateProduct(s testify.Suite, opts ...GenerateProductOption) *domain.Product {
	product := &productToCreate{
		id:     valueobjects.ProductID(uuid.New()),
		name:   valueobjects.ProductName(gofakeit.ProductName()),
		price:  valueobjects.ProductPrice(gofakeit.Price(1.0, 1000.0)),
		status: valueobjects.ProductStatusInit,
	}

	for _, opt := range opts {
		s.Require().NoError(opt(product))
	}

	return domain.NewProduct(product.id, product.name, product.price, product.status)
}

// Методы для генерации Order
type (
	orderToCreate struct {
		id       valueobjects.OrderID
		userID   valueobjects.UserID
		status   domain.OrderStatus
		products []domain.Product
	}

	GenerateOrderOption func(*orderToCreate) error
)

func OrderWithID(id valueobjects.OrderID) GenerateOrderOption {
	return func(order *orderToCreate) error {
		order.id = id
		return nil
	}
}

func OrderWithUserID(userID valueobjects.UserID) GenerateOrderOption {
	return func(order *orderToCreate) error {
		order.userID = userID
		return nil
	}
}

func OrderWithStatus(status domain.OrderStatus) GenerateOrderOption {
	return func(order *orderToCreate) error {
		order.status = status
		return nil
	}
}

func OrderWithProducts(products []domain.Product) GenerateOrderOption {
	return func(order *orderToCreate) error {
		order.products = products
		return nil
	}
}

// GenerateOrder генерирует Order с помощью переданных опций
// Если опции не переданы, генерируется Order со случайными данными:
//   - id       - случайный UUID
//   - userID   - случайный uint64
//   - status   - случайны статус заказа
//   - products - один случайный продукт
func GenerateOrder(s testify.Suite, opts ...GenerateOrderOption) *domain.Order {
	order := &orderToCreate{
		id:       valueobjects.OrderID(uuid.New()),
		userID:   valueobjects.UserID(uuid.New()),
		status:   domain.AllOrderStatuses[gofakeit.IntRange(0, len(domain.AllOrderStatuses)-1)],
		products: []domain.Product{*GenerateProduct(s)},
	}

	for _, opt := range opts {
		s.Require().NoError(opt(order))
	}

	return domain.NewOrder(order.id, order.userID, order.status, order.products)
}
