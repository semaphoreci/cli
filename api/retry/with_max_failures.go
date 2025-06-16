package retry

import "time"

const (
	// InitialBackoff is the initial delay between retries
	InitialBackoff = 100 * time.Millisecond
	// MaxBackoff is the maximum delay between retries
	MaxBackoff = 2 * time.Second
)

// RetryWithMaxFailures executes the given operation with retry logic and exponential backoff
func RetryWithMaxFailures(maxFailures int, operation func() error) error {
	numFailures := 0
	backoff := InitialBackoff

	for {
		err := operation()
		if err == nil {
			return nil
		}

		numFailures++
		if numFailures > maxFailures {
			return err
		}

		// Wait before the next retry with exponential backoff
		time.Sleep(backoff)

		// Double the backoff for next iteration, but don't exceed MaxBackoff
		backoff *= 2
		if backoff > MaxBackoff {
			backoff = MaxBackoff
		}
	}
}
