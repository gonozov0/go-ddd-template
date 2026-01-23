package consumerutils

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"

	"fmt"

	loggerutils "go-ddd-template/pkg/logger/utils"
	"go-ddd-template/pkg/sqs"
)

func ReadAndProcessMessageFromYDB[Event any](
	ctx context.Context,
	topicReader *topicreader.Reader,
	processMessage func(Event),
) error {
	message, err := topicReader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	defer func() {
		if err := topicReader.Commit(ctx, message); err != nil {
			err = fmt.Errorf("failed to commit message: %w", err)
			slog.Error("failed to commit message", loggerutils.ErrAttr(err))
		}
	}()

	data, err := io.ReadAll(message)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	var event Event

	err = json.Unmarshal(data, &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	processMessage(event)

	return nil
}

func ReadAndProcessMessageFromSQS[Event any](
	ctx context.Context,
	reader *sqs.Reader,
	processMessage func(Event) error,
) error {
	messages, err := reader.ReadMessages(ctx)
	if err != nil {
		return fmt.Errorf("failed to read messages: %w", err)
	}

	if len(messages) != 1 {
		// Необходимо настроить reader так, чтобы максимум можно было получить только одно сообщение
		return fmt.Errorf("expected 1 message, got %d", len(messages))
	}

	message := messages[0]
	if message.Body == nil {
		return fmt.Errorf("message body is nil")
	}

	var event Event

	err = json.Unmarshal([]byte(*message.Body), &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err = processMessage(event); err != nil {
		return fmt.Errorf("failed to process message: %w", err)
	}

	if message.ReceiptHandle == nil {
		return fmt.Errorf("message receipt_handle is nil")
	}

	err = reader.DeleteMessage(ctx, *message.ReceiptHandle)
	if err != nil {
		return fmt.Errorf(
			"failed to delete message from queue (receipt_handle: %s): %w",
			*message.ReceiptHandle,
			err,
		)
	}

	return nil
}
