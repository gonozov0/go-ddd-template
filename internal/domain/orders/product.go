package orders

import (
	"errors"
	"fmt"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

var (
	ErrInvalidProduct          = errors.New("invalid product")
	ErrProductAlreadyReserved  = errors.New("product already reserved")
	ErrProductAlreadyPublished = errors.New("product already published")
	ErrProductNotFound         = errors.New("product not found")
	ErrReserveInitedProduct    = errors.New("it is forbidden to reserve product with init status")
)

type Product struct {
	id     valueobjects.ProductID
	name   valueobjects.ProductName
	price  valueobjects.ProductPrice
	status valueobjects.ProductStatus
}

func NewProduct(
	id valueobjects.ProductID,
	name valueobjects.ProductName,
	price valueobjects.ProductPrice,
	status valueobjects.ProductStatus,
) *Product {
	return &Product{
		id:     id,
		name:   name,
		price:  price,
		status: status,
	}
}

func (p *Product) Reserve() error {
	if p.status == valueobjects.ProductStatusInit {
		return fmt.Errorf("%w: %s", ErrReserveInitedProduct, p.id)
	}

	if p.status == valueobjects.ProductStatusReserved {
		return fmt.Errorf("%w: %s", ErrProductAlreadyReserved, p.id)
	}

	p.status = valueobjects.ProductStatusReserved

	return nil
}

func (p *Product) Publish() error {
	if p.status == valueobjects.ProductStatusPublished {
		return fmt.Errorf("%w: %s", ErrProductAlreadyPublished, p.id)
	}

	p.status = valueobjects.ProductStatusPublished

	return nil
}

func (p *Product) GetID() valueobjects.ProductID {
	return p.id
}

func (p *Product) GetName() valueobjects.ProductName {
	return p.name
}

func (p *Product) GetPrice() valueobjects.ProductPrice {
	return p.price
}

func (p *Product) IsPublished() bool {
	return p.status == valueobjects.ProductStatusPublished
}

func (p *Product) GetStatus() valueobjects.ProductStatus {
	return p.status
}

type Products []Product
