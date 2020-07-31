package ch04

import (
	"context"
	"errors"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint64) Circuit {
	// Whether the last interaction with the downstream service was NOT an error
	var lastStateSuccessful = true

	// Number of failures after the first
	var consecutiveFailures uint64 = 0

	// Time of the last interaction with the downstream service
	var lastAttempt time.Time = time.Now()

	// Construct and return the Circuit closure
	return func(ctx context.Context) (string, error) {
		if consecutiveFailures >= failureThreshold {
			// When the circuit should automatically reclose. Note the
			// exponential backoff.
			backoffLevel := consecutiveFailures - failureThreshold
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << backoffLevel)

			if !time.Now().After(shouldRetryAt) {
				return "", errors.New("circuit open -- service unreachable")
			}
		}

		// Call the circuit function and react to any error

		lastAttempt = time.Now()
		response, err := circuit(ctx)
		if err != nil {
			if !lastStateSuccessful {
				consecutiveFailures++
			}
			lastStateSuccessful = false
			return response, err
		}

		lastStateSuccessful = true
		consecutiveFailures = 0

		return response, nil
	}
}
