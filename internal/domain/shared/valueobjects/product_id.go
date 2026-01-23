package valueobjects

import (
	"github.com/google/uuid"

	"errors"
	"fmt"
)

type ProductID uuid.UUID

var ErrInvalidProductID = errors.New("invalid product id")

func NewProductID(id uuid.UUID) (ProductID, error) {
	if id == uuid.Nil {
		return ProductID{}, ErrInvalidProductID
	}

	return ProductID(id), nil
}

func NewProductIDFromString(id string) (ProductID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ProductID{}, ErrInvalidProductID
	}

	return NewProductID(parsedID)
}

func (id ProductID) String() string {
	return uuid.UUID(id).String()
}

func (id ProductID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type ProductIDs []ProductID

func NewProductIDs(ids []uuid.UUID) (ProductIDs, error) {
	productIDs := make(ProductIDs, 0, len(ids))

	for _, id := range ids {
		productID, err := NewProductID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to init ProductID: %w", err)
		}

		productIDs = append(productIDs, productID)
	}

	return productIDs, nil
}

func NewProductIDsFromStrings(ids []string) (ProductIDs, error) {
	productIDs := make(ProductIDs, 0, len(ids))

	for _, id := range ids {
		productID, err := NewProductIDFromString(id)
		if err != nil {
			return nil, fmt.Errorf("failed to init ProductID: %w", err)
		}

		productIDs = append(productIDs, productID)
	}

	return productIDs, nil
}

func (ids ProductIDs) UUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uuids = append(uuids, id.UUID())
	}

	return uuids
}

func (ids ProductIDs) Strings() []string {
	strings := make([]string, 0, len(ids))
	for _, id := range ids {
		strings = append(strings, id.String())
	}

	return strings
}
