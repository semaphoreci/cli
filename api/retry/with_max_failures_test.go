package retry

import (
	"errors"
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
