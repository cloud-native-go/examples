package ch04

import (
	"context"
	"sync"
	"time"
)

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)
		threshold = time.Now().Add(d)

		return result, err
	}
}
