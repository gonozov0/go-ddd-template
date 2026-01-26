package valueobjects

import (
	"fmt"
)

type ProductStatus string

const (
	ProductStatusInit      ProductStatus = "init"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusReserved  ProductStatus = "reserved"
)

func NewProductStatus(status string) (ProductStatus, error) {
	switch ProductStatus(status) {
	case ProductStatusInit:
		return ProductStatusInit, nil
	case ProductStatusPublished:
		return ProductStatusPublished, nil
	case ProductStatusReserved:
		return ProductStatusReserved, nil
	default:
		return "", fmt.Errorf("invalid product status: %s", status)
	}
}

func (s ProductStatus) String() string {
	return string(s)
}
