package consumerutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/traces"
)

type MessageHandler func(ctx context.Context, data []byte) error

//nolint:cyclop
func RunYDBConsumer(
	ctx context.Context,
	reader *topicreader.Reader,
	name string,
	handler MessageHandler,
) error {
	var (
		msg *topicreader.Message
		err error
	)

	handler = panicsHandlerMiddleware(handler)
	handler = tracesHandlerMiddleware(handler, name)
	handler = contextCancelMiddleware(handler)

	for {
		if msg, err = reader.ReadMessage(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("failed to read message: %w", err)
		}

		bo := backoff.NewExponentialBackOff()
		bo.MaxElapsedTime = 0

		data, err := io.ReadAll(msg)
		if err != nil {
			return fmt.Errorf("failed to read message body: %w", err)
		}

		if err = backoff.RetryNotify(
			func() error { return handler(ctx, data) },
			backoff.WithContext(bo, ctx),
			logRetry(msg),
		); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("failed to handle message: %w", err)
		}

		if err = reader.Commit(context.Background(), msg); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}
	}
}

func contextCancelMiddleware(h MessageHandler) MessageHandler {
	return func(ctx context.Context, data []byte) error {
		if ctx.Err() != nil {
			//nolint:descriptiveerrors
			return backoff.Permanent(ctx.Err())
		}

		//nolint:descriptiveerrors
		return h(context.Background(), data)
	}
}

func panicsHandlerMiddleware(h MessageHandler) MessageHandler {
	return func(ctx context.Context, data []byte) (err error) {
		defer func() {
			if v := recover(); v != nil {
				err = fmt.Errorf("panic during message handling: %v", v)
			}
		}()

		return h(ctx, data)
	}
}

func tracesHandlerMiddleware(h MessageHandler, name string) MessageHandler {
	return func(ctx context.Context, data []byte) error {
		correlationID := uuid.New().String()

		ctx = traces.WithCorrelationID(ctx, correlationID)

		ctx, span := traces.CreateSpan(
			ctx,
			traces.ScopeConsumer,
			fmt.Sprintf("%s %s", traces.ScopeTraceMiddleware, name),
			trace.SpanKindConsumer,
		)
		defer span.End()

		span.SetAttributes(
			attribute.String("correlation_id", correlationID),
		)

		err := h(ctx, data)

		traces.SetSpanStatus(span, err)

		return err
	}
}

func logRetry(msg *topicreader.Message) func(error, time.Duration) {
	return func(err error, d time.Duration) {
		slog.Error("retrying message processing",
			"backoff_in", d,
			"offset", msg.Offset,
			"seq", msg.SeqNo,
			loggerutils.ErrAttr(err),
		)
	}
}
