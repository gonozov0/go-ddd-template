package consumers_test

import (
	"context"

	"github.com/google/uuid"

	deliveryunithelpers "go-ddd-template/internal/application/deliveries/shared/helpers"
)

func (s *DeliveriesSuite) TestHandelOrderCreated() {
	s.Run("successful create delivery from order created", func() {
		startCount := s.getDeliveriesCount(context.Background())

		orderID := uuid.New()
		payload := deliveryunithelpers.GenerateOrderCreatedEventPayload(s,
			deliveryunithelpers.OrderCreatedEventWithOrderID(orderID),
		)

		err := s.DeliveryConsumers.HandleOrderProcessing(context.Background(), payload)
		s.Require().NoError(err)

		deliveries, err := s.DeliveriesRepo.ListDeliveries(context.Background())
		s.Require().NoError(err)
		s.Require().Len(deliveries, startCount+1)

		s.Require().Equal(orderID.String(), deliveries[0].GetOrderID().String())
	})
	s.Run("creates multiple deliveries for multiple events", func() {
		startCount := s.getDeliveriesCount(context.Background())

		err := s.DeliveryConsumers.HandleOrderProcessing(
			context.Background(),
			deliveryunithelpers.GenerateOrderCreatedEventPayload(s),
		)
		s.Require().NoError(err)

		err = s.DeliveryConsumers.HandleOrderProcessing(
			context.Background(),
			deliveryunithelpers.GenerateOrderCreatedEventPayload(s),
		)
		s.Require().NoError(err)

		deliveries, err := s.DeliveriesRepo.ListDeliveries(context.Background())
		s.Require().NoError(err)
		s.Require().Len(deliveries, startCount+2)
	})
	s.Run("fails with invalid json payload", func() {
		payload := []byte("invalid json")

		err := s.DeliveryConsumers.HandleOrderProcessing(context.Background(), payload)
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "failed to unmarshal")
	})
	s.Run("fails with invalid order uuid format", func() {
		payload := []byte(`{"order_id": "not-a-uuid"}`)

		err := s.DeliveryConsumers.HandleOrderProcessing(context.Background(), payload)
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "failed to unmarshal")
	})
	s.Run("fails with nil uuid", func() {
		startCount := s.getDeliveriesCount(context.Background())

		payload := deliveryunithelpers.GenerateOrderCreatedEventPayload(s,
			deliveryunithelpers.OrderCreatedEventWithOrderID(uuid.Nil),
		)

		err := s.DeliveryConsumers.HandleOrderProcessing(context.Background(), payload)

		s.Require().Error(err)
		s.Require().Contains(err.Error(), "failed to create order id")

		deliveries, err := s.DeliveriesRepo.ListDeliveries(context.Background())
		s.Require().NoError(err)
		s.Require().Len(deliveries, startCount)
	})
	s.Run("handles context cancellation gracefully", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		payload := deliveryunithelpers.GenerateOrderCreatedEventPayload(s)

		err := s.DeliveryConsumers.HandleOrderProcessing(ctx, payload)
		if err != nil {
			s.Require().ErrorIs(err, context.Canceled)
		}
	})
}
