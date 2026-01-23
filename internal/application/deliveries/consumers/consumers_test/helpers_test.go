package consumers_test

import (
	"context"
)

func (s *DeliveriesSuite) getDeliveriesCount(ctx context.Context) int {
	deliveries, err := s.DeliveriesRepo.ListDeliveries(ctx)
	s.Require().NoError(err)

	return len(deliveries)
}
