package helpers

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/testify"
)

// Методы для генерации Product
type (
	productToCreate struct {
		id            valueobjects.ProductID
		name          valueobjects.ProductName
		price         valueobjects.ProductPrice
		status        valueobjects.ProductStatus
		imageFilename valueobjects.ImageFilename
		imageURL      valueobjects.ImageURL
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

func ProductWithImageFilename(imageFilename valueobjects.ImageFilename) GenerateProductOption {
	return func(product *productToCreate) error {
		product.imageFilename = imageFilename
		return nil
	}
}

func ProductWithImageURL(imageURL valueobjects.ImageURL) GenerateProductOption {
	return func(product *productToCreate) error {
		product.imageURL = imageURL
		return nil
	}
}

// GenerateProduct генерирует Product с помощью переданных опций
// Если опции не переданы, генерируется Product со случайными данными:
//   - id		- случайный UUID
//   - name		- случайное название продукта
//   - price	- случайная цена от 1.0 до 1000.0
//   - status   - статус Init
func GenerateProduct(s testify.Suite, opts ...GenerateProductOption) *domain.Product {
	product := &productToCreate{
		id:            valueobjects.ProductID(uuid.New()),
		name:          valueobjects.ProductName(gofakeit.ProductName()),
		price:         valueobjects.ProductPrice(gofakeit.Price(1.0, 1000.0)),
		status:        valueobjects.ProductStatusInit,
		imageFilename: valueobjects.EmptyImageFilename,
		imageURL:      valueobjects.EmptyImageURL,
	}

	for _, opt := range opts {
		s.Require().NoError(opt(product))
	}

	return domain.NewProduct(
		product.id,
		product.name,
		product.price,
		product.status,
		product.imageFilename,
		product.imageURL,
	)
}
