package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	reciveMessageOptions struct {
		WaitTimeSeconds     int64
		MaxNumberOfMessages int64
		VisabilityTimeout   int64
	}

	ReciveMessageOption func(*reciveMessageOptions)
)

func WithWaitTimeSeconds(waitTimeSeconds int64) ReciveMessageOption {
	return func(opts *reciveMessageOptions) {
		opts.WaitTimeSeconds = waitTimeSeconds
	}
}

func WithMaxNumberOfMessages(maxNumberOfMessages int64) ReciveMessageOption {
	return func(opts *reciveMessageOptions) {
		opts.MaxNumberOfMessages = maxNumberOfMessages
	}
}

func WithVisabilityTimeout(visabilityTimeout int64) ReciveMessageOption {
	return func(opts *reciveMessageOptions) {
		opts.VisabilityTimeout = visabilityTimeout
	}
}

type Reader struct {
	sqs                  *sqs.SQS
	queueURL             string
	reciveMessageOptions reciveMessageOptions
}

func NewReader(sqsClient Session, queueName string, opts ...ReciveMessageOption) *Reader {
	reader := Reader{
		sqs:      sqsClient.SQS,
		queueURL: fmt.Sprintf("%s_%s", sqsClient.QueueBaseName, queueName),
		reciveMessageOptions: reciveMessageOptions{
			WaitTimeSeconds:     20,
			MaxNumberOfMessages: 10,
			VisabilityTimeout:   60,
		},
	}

	for _, opt := range opts {
		opt(&reader.reciveMessageOptions)
	}

	return &reader
}

func (r *Reader) ReadMessages(ctx context.Context) ([]*sqs.Message, error) {
	receiveMessageOutput, err := r.sqs.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		WaitTimeSeconds: aws.Int64(
			r.reciveMessageOptions.WaitTimeSeconds,
		), // время ожидания получения хотя бы одного сообщения
		MaxNumberOfMessages: aws.Int64(
			r.reciveMessageOptions.MaxNumberOfMessages,
		), // максимальное количество сообщений, которое можно получить
		VisibilityTimeout: aws.Int64(
			r.reciveMessageOptions.VisabilityTimeout,
		), // время, в течение которого сообщения не будут видны для других потребителей
		QueueUrl: aws.String(r.queueURL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from queue (queueUrl: %s): %w", r.queueURL, err)
	}

	return receiveMessageOutput.Messages, nil
}

func (r *Reader) MakeMessageVisible(ctx context.Context, receiptHandle string) error {
	_, err := r.sqs.ChangeMessageVisibilityWithContext(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:      aws.String(r.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
		VisibilityTimeout: aws.Int64(
			1,
		), // Сообщение появится в очереди через одну секунду. Требуется для синхронизации. Возникали ситуации, когда сообщение появлялось в очереди, а транзакция, внутри которой вызывалось это действие, не успевала закомититься
	})
	if err != nil {
		return fmt.Errorf(
			"failed to return message to queue (queueUrl: %s, receiptHandle: %s): %w",
			r.queueURL,
			receiptHandle,
			err,
		)
	}

	return nil
}

func (r *Reader) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := r.sqs.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(r.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf(
			"failed to delete message from queue (queueUrl: %s, receiptHandle: %s): %w",
			r.queueURL,
			receiptHandle,
			err,
		)
	}

	return nil
}
