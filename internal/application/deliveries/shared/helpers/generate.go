package helpers

import (
	"encoding/json"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/testify"
)

// Методы для генерации Delivery
type (
	deliveryToCreate struct {
		id        valueobjects.DeliveryID
		orderID   valueobjects.OrderID
		createdAt valueobjects.Timestamp
	}

	GenerateDeliveryOption func(*deliveryToCreate) error
)

func DeliveryWithID(id valueobjects.DeliveryID) GenerateDeliveryOption {
	return func(delivery *deliveryToCreate) error {
		delivery.id = id
		return nil
	}
}

func DeliveryWithOrderID(orderID valueobjects.OrderID) GenerateDeliveryOption {
	return func(delivery *deliveryToCreate) error {
		delivery.orderID = orderID
		return nil
	}
}

func DeliveryWithCreatedAt(createdAt valueobjects.Timestamp) GenerateDeliveryOption {
	return func(delivery *deliveryToCreate) error {
		delivery.createdAt = createdAt
		return nil
	}
}

// GenerateDelivery генерирует Delivery с помощью переданных опций
// Если опции не переданы, генерируется Delivery со случайными данными:
//   - id        - случайный UUID
//   - orderID   - случайный UUID
//   - createdAt - случайная дата
func GenerateDelivery(s testify.Suite, opts ...GenerateDeliveryOption) *domain.Delivery {
	deliveryID, _ := valueobjects.NewDeliveryID(uuid.New())
	orderID, _ := valueobjects.NewOrderID(uuid.New())

	delivery := &deliveryToCreate{
		id:        deliveryID,
		orderID:   orderID,
		createdAt: valueobjects.NewTimestamp(gofakeit.Date()),
	}

	for _, opt := range opts {
		s.Require().NoError(opt(delivery))
	}

	return domain.NewDelivery(delivery.id, delivery.orderID, delivery.createdAt)
}

// Методы для генерации события orderCreated
type (
	orderCreatedEventToCreate struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	OrderCreatedEventOption func(*orderCreatedEventToCreate) error
)

func OrderCreatedEventWithOrderID(orderID uuid.UUID) OrderCreatedEventOption {
	return func(orderCreatedEvent *orderCreatedEventToCreate) error {
		orderCreatedEvent.OrderID = orderID
		return nil
	}
}

// GenerateOrderCreatedEventPayload генерирует orderCreatedEvent с помощью переданных опций
// Если опции не переданы, генерируется orderCreatedEvent со случайными данными:
//   - orderID - случайный UUID
func GenerateOrderCreatedEventPayload(s testify.Suite, opts ...OrderCreatedEventOption) []byte {
	orderCreatedEventToCreate := &orderCreatedEventToCreate{
		OrderID: uuid.New(),
	}

	for _, opt := range opts {
		s.Require().NoError(opt(orderCreatedEventToCreate))
	}

	payload, _ := json.Marshal(&orderCreatedEventToCreate)

	return payload
}
