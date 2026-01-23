package backoff

import (
	"context"
	"log/slog"
	"net"

	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

const DefaultRetryLimit = 3

type RetryableError interface {
	error
	IsRetryable() bool
}

func RunWithRetry(ctx context.Context, fn func() error, retryLimit uint64) error {
	return backoff.Retry(func() error {
		err := fn()
		if err != nil {
			if isRetryable(err) {
				slog.WarnContext(
					ctx,
					"Backoff triggered",
					loggerutils.ErrAttr(fmt.Errorf("backoff trigger warning: %w", err)),
				)
				//nolint:descriptiveerrors
				return err
			}

			//nolint:descriptiveerrors
			return backoff.Permanent(err)
		}

		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), retryLimit))
}

func isRetryable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return true
	}

	var retriableErr RetryableError
	if errors.As(err, &retriableErr) {
		return retriableErr.IsRetryable()
	}

	return false
}
