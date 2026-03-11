package retry

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	// InitialBackoff is the initial delay between retries
	InitialBackoff = 100 * time.Millisecond
	// MaxBackoff is the maximum delay between retries
	MaxBackoff = 2 * time.Second
)

// NonRetryableError wraps an error that should not be retried.
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return e.Err.Error()
}

func (e *NonRetryableError) Unwrap() error {
	return e.Err
}

// NonRetryable marks an error as non-retryable.
func NonRetryable(err error) error {
	return &NonRetryableError{Err: err}
}

// RetryWithMaxFailures executes the given operation with retry logic and exponential backoff.
// Non-retryable errors (wrapped with NonRetryable) are returned immediately.
func RetryWithMaxFailures(maxFailures int, operation func() error) error {
	numFailures := 0
	backoff := InitialBackoff

	for {
		err := operation()
		if err == nil {
			return nil
		}

		var nonRetryable *NonRetryableError
		if errors.As(err, &nonRetryable) {
			return nonRetryable.Err
		}

		numFailures++
		if numFailures > maxFailures {
			return fmt.Errorf("failed after %d attempts: %w", numFailures, err)
		}

		log.Printf("attempt %d/%d failed (%v), retrying in %v", numFailures, maxFailures, err, backoff)

		// Wait before the next retry with exponential backoff
		time.Sleep(backoff)

		// Double the backoff for next iteration, but don't exceed MaxBackoff
		backoff *= 2
		if backoff > MaxBackoff {
			backoff = MaxBackoff
		}
	}
}
