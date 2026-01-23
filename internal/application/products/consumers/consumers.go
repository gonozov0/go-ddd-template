package consumers

import (
	"context"

	"golang.org/x/sync/errgroup"

	service "go-ddd-template/internal/service/products"
	"go-ddd-template/internal/shared"
	"go-ddd-template/pkg/consumerutils"
	"go-ddd-template/pkg/sqs"
)

type QueueReaders struct {
	ProductInited *sqs.Reader
}

type DeliveryConsumers struct {
	productService service.ProductService
}

func NewDeliveryConsumers(repo service.ProductRepository, imageStorage service.ImageStorage) *DeliveryConsumers {
	return &DeliveryConsumers{
		productService: service.NewProductService(repo, imageStorage),
	}
}

func (d DeliveryConsumers) Start(
	ctx context.Context,
	g *errgroup.Group,
	queueReaders QueueReaders,
) {
	g.Go(func() error {
		return consumerutils.RunSQSConsumer(
			ctx,
			queueReaders.ProductInited,
			shared.ProductInitedConsumer,
			d.CheckProduct,
		)
	})
}
