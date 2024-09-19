package create

import (
	"context"
	"errors"

	"go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type ProductsAlreadyReservedError struct {
	ProductIDs []uuid.UUID
}

func (e ProductsAlreadyReservedError) Error() string {
	return "products already reserved"
}

func (s *OrderCreationService) reserveProducts(ctx context.Context, items []Item) ([]products.Product, error) {
	ps, err := s.productRepo.GetProductsForUpdate(ctx, getItemIDs(items))
	if err != nil {
		return nil, err
	}

	var reservedProductIDs []uuid.UUID
	for _, product := range ps {
		if err := product.Reserve(); err != nil {
			if errors.Is(err, products.ErrProductAlreadyReserved) {
				reservedProductIDs = append(reservedProductIDs, product.ID())
				continue
			}
			return nil, err
		}
	}
	if len(reservedProductIDs) > 0 {
		return nil, ProductsAlreadyReservedError{ProductIDs: reservedProductIDs}
	}

	return ps, nil
}

func getItemIDs(items []Item) []uuid.UUID {
	itemIDs := make([]uuid.UUID, 0, len(items))
	for _, i := range items {
		itemIDs = append(itemIDs, i.ID)
	}
	return itemIDs
}
