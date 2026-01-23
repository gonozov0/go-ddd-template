package deliveries

import "errors"

var (
	ErrDeliveryNotFound = errors.New("delivery not found")
)

type RetriableDeliveryNotFoundError struct {
	err error
}

func NewRetriableDeliveryNotFoundError() *RetriableDeliveryNotFoundError {
	return &RetriableDeliveryNotFoundError{err: ErrDeliveryNotFound}
}

func (e *RetriableDeliveryNotFoundError) Error() string {
	return e.err.Error()
}

func (e *RetriableDeliveryNotFoundError) IsRetryable() bool {
	return true
}

func (e *RetriableDeliveryNotFoundError) Unwrap() error {
	return e.err
}
