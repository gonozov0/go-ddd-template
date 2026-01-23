package valueobjects

import (
	"errors"
	"fmt"
)

type ProductPrice float64

var ErrInvalidProductPrice = errors.New("invalid product price")

func NewProductPrice(price float64) (ProductPrice, error) {
	if price <= 0 {
		return 0, fmt.Errorf("%w: price must be greater than 0", ErrInvalidProductPrice)
	}

	return ProductPrice(price), nil
}

func (p ProductPrice) Float64() float64 {
	return float64(p)
}
