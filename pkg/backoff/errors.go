package backoff

import (
	"fmt"
	"net/http"
	"slices"
)

type (
	HTTPError struct {
		StatusCode int
	}

	AlwaysRetryableError struct {
		err error
	}
)

func NewHTTPError(statusCode int) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

var (
	defaultStatusCodesForBackoff = []int{
		http.StatusTooManyRequests,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}
)

func CheckStatusCodeForBackoff(actualStatusCode int, expectedStatusCodesForBackoff ...int) error {
	if len(expectedStatusCodesForBackoff) == 0 {
		expectedStatusCodesForBackoff = defaultStatusCodesForBackoff
	}

	if slices.Contains(expectedStatusCodesForBackoff, actualStatusCode) {
		return NewHTTPError(actualStatusCode)
	}

	return nil
}

func NewAlwaysRetryableError(err error) *AlwaysRetryableError {
	return &AlwaysRetryableError{err: err}
}

func (e *AlwaysRetryableError) Error() string {
	return e.err.Error()
}

func (e *AlwaysRetryableError) IsRetryable() bool {
	return true
}

func (e *AlwaysRetryableError) Unwrap() error {
	return e.err
}
