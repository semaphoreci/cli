package retry

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRetryWithMaxFailures(t *testing.T) {
	tests := []struct {
		name          string
		maxFailures   int
		failureCount  int
		expectedError bool
		expectedCalls int
		operationTime time.Duration
	}{
		{
			name:          "success on first try",
			maxFailures:   3,
			failureCount:  0,
			expectedError: false,
			expectedCalls: 1,
		},
		{
			name:          "success after one failure",
			maxFailures:   3,
			failureCount:  1,
			expectedError: false,
			expectedCalls: 2,
		},
		{
			name:          "exceeds max failures",
			maxFailures:   2,
			failureCount:  3,
			expectedError: true,
			expectedCalls: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			start := time.Now()

			err := RetryWithMaxFailures(tt.maxFailures, func() error {
				callCount++
				if callCount <= tt.failureCount {
					return errors.New("test error")
				}
				return nil
			})

			duration := time.Since(start)

			if tt.expectedError && err == nil {
				t.Error("expected an error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if callCount != tt.expectedCalls {
				t.Errorf("expected %d calls but got %d", tt.expectedCalls, callCount)
			}

			// Verify that backoff is working
			if tt.failureCount > 0 {
				// We should have at least waited InitialBackoff
				if duration < InitialBackoff {
					t.Errorf("expected duration >= %v, got %v", InitialBackoff, duration)
				}
			}
		})
	}
}

func TestRetryWithMaxFailuresBackoff(t *testing.T) {
	maxFailures := 3
	callCount := 0
	var delays []time.Duration
	lastCall := time.Now()

	err := RetryWithMaxFailures(maxFailures, func() error {
		if callCount > 0 {
			delay := time.Since(lastCall)
			delays = append(delays, delay)
		}
		lastCall = time.Now()
		callCount++
		return errors.New("always fail")
	})

	if err == nil {
		t.Error("expected error but got none")
	}

	if len(delays) != maxFailures {
		t.Errorf("expected %d delays but got %d", maxFailures, len(delays))
	}

	// Verify exponential backoff
	for i := 1; i < len(delays); i++ {
		if delays[i] < delays[i-1] {
			t.Errorf("expected increasing delays, but delay %d (%v) < delay %d (%v)",
				i, delays[i], i-1, delays[i-1])
		}
	}
}

func TestRetryWithMaxFailures_NonRetryable(t *testing.T) {
	callCount := 0

	err := RetryWithMaxFailures(5, func() error {
		callCount++
		return NonRetryable(errors.New("bad request"))
	})

	if err == nil {
		t.Fatal("expected an error but got none")
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retries) but got %d", callCount)
	}
	if err.Error() != "bad request" {
		t.Errorf("expected unwrapped error message, got: %v", err)
	}

	// Should not be wrapped in NonRetryableError after return
	var nonRetryable *NonRetryableError
	if errors.As(err, &nonRetryable) {
		t.Error("returned error should not be wrapped in NonRetryableError")
	}
}

func TestRetryWithMaxFailures_ErrorWrapsAttemptCount(t *testing.T) {
	err := RetryWithMaxFailures(2, func() error {
		return errors.New("connection refused")
	})

	if err == nil {
		t.Fatal("expected an error but got none")
	}
	if !strings.Contains(err.Error(), "failed after 3 attempts") {
		t.Errorf("expected error to contain attempt count, got: %v", err)
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("expected error to contain original message, got: %v", err)
	}

	// Original error should be unwrappable
	if !errors.Is(err, errors.Unwrap(err)) {
		// Just verify Unwrap works
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			t.Error("expected error to be unwrappable")
		}
	}
}

func TestRetryWithMaxFailures_MixedRetryableAndNon(t *testing.T) {
	callCount := 0

	err := RetryWithMaxFailures(5, func() error {
		callCount++
		if callCount <= 2 {
			return errors.New("transient error")
		}
		return NonRetryable(errors.New("permanent error"))
	})

	if err == nil {
		t.Fatal("expected an error but got none")
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls (2 retries then non-retryable) but got %d", callCount)
	}
	if err.Error() != "permanent error" {
		t.Errorf("expected 'permanent error', got: %v", err)
	}
}
