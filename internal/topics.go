package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicoptions"

	deliveriesapp "go-ddd-template/internal/application/deliveries/consumers"
	productsapp "go-ddd-template/internal/application/products/consumers"
	ordersinfra "go-ddd-template/internal/infrastructure/orders"
	productsinfra "go-ddd-template/internal/infrastructure/products"
	"go-ddd-template/internal/shared"
	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

func InitQueueReaders(
	driver *ydb.Driver,
	sqsSession sqs.Session,
) (deliveryQueueReaders deliveriesapp.QueueReaders, productQueueReaders productsapp.QueueReaders, closer func(context.Context), err error) {
	var (
		closers  []func(context.Context)
		topicErr error
	)

	closer = func(ctx context.Context) {
		for _, c := range closers {
			c(ctx)
		}
	}

	defer func(ctx context.Context) {
		if err != nil {
			closer(ctx)
			closer = func(context.Context) {}
		}
	}(context.Background())

	deliveryQueueReaders.OrderProcessingDelivery, topicErr = driver.Topic().StartReader(
		shared.OrderProcessingConsumer,
		topicoptions.ReadTopic(shared.OrderProcessingTopic),
	)
	if topicErr != nil {
		return deliveryQueueReaders, productQueueReaders, closer,
			errors.Join(
				err,
				fmt.Errorf("failed to start %s topic reader: %w", shared.OrderProcessingTopic, topicErr),
			)
	}

	slog.Info("consumer has successfully started", "topic", shared.OrderProcessingTopic)

	productQueueReaders.ProductInited = sqs.NewReader(sqsSession, shared.ProductInitedQueue)
	slog.Info("producer has successfully started", "topic", shared.ProductInitedQueue)

	closers = append(closers, func(ctx context.Context) {
		if closeErr := deliveryQueueReaders.OrderProcessingDelivery.Close(ctx); closeErr != nil {
			slog.Error(
				"failed to close consumer",
				"topic",
				shared.OrderProcessingTopic,
				loggerutils.ErrAttr(closeErr),
			)

			return
		}

		slog.Info("consumer stopped gracefully", "topic", shared.OrderProcessingTopic)
	})

	return deliveryQueueReaders, productQueueReaders, closer, nil
}

func InitQueueWriters(
	driver *ydb.Driver,
	sqsSession sqs.Session,
) (orderQueueWriters ordersinfra.QueueWriters, productQueueWriters productsinfra.QueueWriters, closer func(context.Context), err error) {
	var (
		closers  []func(context.Context)
		topicErr error
	)

	closer = func(ctx context.Context) {
		for _, c := range closers {
			c(ctx)
		}
	}

	defer func(ctx context.Context) {
		if err != nil { // если при инициализации writer'ов возникла ошибка, то сразу закрываем соединения
			closer(ctx)
			closer = func(context.Context) {} // определяю как пустую функцию, чтобы никто случайно дважды не вызвал функцию close
		}
	}(context.Background())

	orderQueueWriters.OrderCreated, topicErr = driver.Topic().StartWriter(shared.OrderProcessingTopic)
	if topicErr != nil {
		return orderQueueWriters, productQueueWriters, closer,
			errors.Join(
				err,
				fmt.Errorf("failed to start %s topic writer: %w", shared.OrderProcessingTopic, topicErr),
			)
	}

	slog.Info("producer has successfully started", "topic", shared.OrderProcessingTopic)

	productQueueWriters.ProductInited = sqs.NewWriter(sqsSession, shared.ProductInitedQueue)
	slog.Info("producer has successfully started", "topic", shared.ProductInitedQueue)

	closers = append(closers, func(ctx context.Context) {
		if closeErr := orderQueueWriters.OrderCreated.Close(ctx); closeErr != nil {
			slog.Error(
				"failed to close producer",
				"topic",
				shared.OrderProcessingTopic,
				loggerutils.ErrAttr(closeErr),
			)

			return
		}

		slog.Info("producer stopped gracefully", "topic", shared.OrderProcessingTopic)
	})

	return orderQueueWriters, productQueueWriters, closer, nil
}
