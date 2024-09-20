package orders

import (
	"errors"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrProductAlreadyReserved = errors.New("product already reserved")
	ErrProductNotFound        = errors.New("product not found")
)
