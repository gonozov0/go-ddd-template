package consumerutils

import (
	"context"
	"log/slog"

	loggerutils "go-ddd-template/pkg/logger/utils"
	libsqs "go-ddd-template/pkg/sqs"
)

//nolint:cyclop
func RunSQSConsumer(
	ctx context.Context,
	reader *libsqs.Reader,
	name string,
	handler MessageHandler,
) error {
	handler = panicsHandlerMiddleware(handler)
	handler = tracesHandlerMiddleware(handler, name)
	handler = contextCancelMiddleware(handler)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			ctx_ := context.Background()

			messages, err := reader.ReadMessages(ctx_)
			if err != nil {
				slog.Error("failed to receive messages", loggerutils.ErrAttr(err))
				continue
			}

			for _, message := range messages {
				if message.MessageId == nil {
					slog.Error("recieve message with nil message id")
					continue
				}

				if message.Body == nil {
					slog.Error("recieve message with nil body")
					continue
				}

				if err = handler(ctx_, []byte(*message.Body)); err != nil {
					slog.Error("failed to handle message", loggerutils.ErrAttr(err))

					if err := reader.MakeMessageVisible(ctx_, *message.ReceiptHandle); err != nil {
						slog.Error("failed to return message to queue", loggerutils.ErrAttr(err))
					}

					continue
				}

				if err = reader.DeleteMessage(ctx_, *message.ReceiptHandle); err != nil {
					slog.Error("failed to delete message", loggerutils.ErrAttr(err))
					continue
				}
			}
		}
	}
}
