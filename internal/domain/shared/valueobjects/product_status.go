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

func (s ProductStatus) String() string {
	return string(s)
}

func NewProductStatus(status string) (ProductStatus, error) {
	switch status {
	case ProductStatusInit.String():
		return ProductStatusInit, nil
	case ProductStatusPublished.String():
		return ProductStatusPublished, nil
	case ProductStatusReserved.String():
		return ProductStatusReserved, nil
	default:
		return "", fmt.Errorf("invalid product status: %s", status)
	}
}
