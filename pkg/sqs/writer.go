package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	sendMessageOptions struct {
		DelaySeconds int64
	}

	SendMessageOption func(*sendMessageOptions)
)

func WithDelaySeconds(delaySeconds int64) SendMessageOption {
	return func(opts *sendMessageOptions) {
		opts.DelaySeconds = delaySeconds
	}
}

type Writer struct {
	sqs                *sqs.SQS
	queueURL           string
	sendMessageOptions sendMessageOptions
}

func NewWriter(sqsClient Session, queueName string, opts ...SendMessageOption) *Writer {
	writer := Writer{
		sqs:      sqsClient.SQS,
		queueURL: fmt.Sprintf("%s_%s", sqsClient.QueueBaseName, queueName),

		sendMessageOptions: sendMessageOptions{
			DelaySeconds: 1,
		},
	}

	for _, opt := range opts {
		opt(&writer.sendMessageOptions)
	}

	return &writer
}

func (p *Writer) SendMessage(ctx context.Context, message string) error {
	sendMessageInput := sqs.SendMessageInput{
		QueueUrl:     aws.String(p.queueURL),
		MessageBody:  aws.String(message),
		DelaySeconds: aws.Int64(p.sendMessageOptions.DelaySeconds),
	}

	_, err := p.sqs.SendMessageWithContext(ctx, &sendMessageInput)
	if err != nil {
		return fmt.Errorf("failed to send message to queue (queueUrl: %s): %w", p.queueURL, err)
	}

	return nil
}
