package valueobjects

import (
	"github.com/google/uuid"

	"errors"
)

type DeliveryID uuid.UUID

var (
	ErrInvalidDeliveryID = errors.New("invalid delivery id")
	EmptyDeliveryID      = DeliveryID{}
)

func NewDeliveryID(id uuid.UUID) (DeliveryID, error) {
	if id == uuid.Nil {
		return DeliveryID{}, ErrInvalidDeliveryID
	}

	return DeliveryID(id), nil
}

func NewDeliveryIDFromString(id string) (DeliveryID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return DeliveryID{}, ErrInvalidDeliveryID
	}

	return NewDeliveryID(parsedID)
}

func (id DeliveryID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id DeliveryID) String() string {
	return uuid.UUID(id).String()
}

func (id DeliveryID) IsEmpty() bool {
	return id == EmptyDeliveryID
}
