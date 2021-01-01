package ch04

import (
	"context"
	"time"
)

func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time = time.Now()
	var ticker *time.Ticker
	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		threshold = time.Now().Add(d)

		if ticker == nil {
			ticker = time.NewTicker(time.Millisecond * 100)
			tickerc := ticker.C

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-tickerc:
						if threshold.Before(time.Now()) {
							result, err = circuit(ctx)
							ticker.Stop()
							ticker = nil
							break
						}
					case <-ctx.Done():
						result, err = "", ctx.Err()
						break
					}
				}
			}()
		}

		return result, err
	}
}
