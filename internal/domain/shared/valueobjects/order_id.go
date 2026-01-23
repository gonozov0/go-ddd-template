package valueobjects

import (
	"github.com/google/uuid"

	"errors"
)

type OrderID uuid.UUID

var EmptyOrderID = OrderID(uuid.Nil)

var ErrInvalidOrderID = errors.New("invalid order id")

func NewOrderID(id uuid.UUID) (OrderID, error) {
	if id == uuid.Nil {
		return EmptyOrderID, ErrInvalidOrderID
	}

	return OrderID(id), nil
}

func NewOrderIDFromString(id string) (OrderID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return EmptyOrderID, ErrInvalidOrderID
	}

	return NewOrderID(parsedID)
}

func (id OrderID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id OrderID) String() string {
	return uuid.UUID(id).String()
}

type OrderIDs []OrderID

func (ids OrderIDs) UUIDs() uuid.UUIDs {
	uuids := make(uuid.UUIDs, 0, len(ids))
	for _, id := range ids {
		uuids = append(uuids, id.UUID())
	}

	return uuids
}
