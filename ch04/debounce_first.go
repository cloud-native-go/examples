package ch04

import (
	"context"
	"time"
)

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var cResult string
	var cError error

	return func(ctx context.Context) (string, error) {
		if threshold.Before(time.Now()) {
			cResult, cError = circuit(ctx)
		}

		threshold = time.Now().Add(d)
		return cResult, cError
	}
}
