package consumers

import (
	"context"

	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
	"golang.org/x/sync/errgroup"

	service "go-ddd-template/internal/service/deliveries"
	"go-ddd-template/internal/shared"
	"go-ddd-template/pkg/consumerutils"
)

type QueueReaders struct {
	OrderProcessingDelivery *topicreader.Reader
}

type DeliveryConsumers struct {
	deliveryService service.DeliveryService
}

func NewDeliveryConsumers(repo service.DeliveryRepository) *DeliveryConsumers {
	return &DeliveryConsumers{
		deliveryService: service.NewDeliveryService(repo),
	}
}

func (d DeliveryConsumers) Start(
	ctx context.Context,
	g *errgroup.Group,
	queueReaders QueueReaders,
) {
	g.Go(func() error {
		return consumerutils.RunYDBConsumer(
			ctx,
			queueReaders.OrderProcessingDelivery,
			shared.OrderProcessingConsumer,
			d.HandleOrderProcessing,
		)
	})
}
