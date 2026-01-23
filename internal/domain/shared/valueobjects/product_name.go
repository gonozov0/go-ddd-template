package valueobjects

import (
	"errors"
	"fmt"
)

type ProductName string

var ErrInvalidProductName = errors.New("invalid product name")

func NewProductName(name string) (ProductName, error) {
	if name == "" {
		return "", fmt.Errorf("%w: name must not be empty", ErrInvalidProductName)
	}

	return ProductName(name), nil
}

func (n ProductName) String() string {
	return string(n)
}
