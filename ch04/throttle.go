package ch04

import (
	"context"
	"time"
)

type Effector func(context.Context) (string, error)

func Throttle(effector Effector, max uint, refill uint, duration time.Duration) Effector {
	var ticker *time.Ticker = nil
	var tokens uint = max

	var lastReturnString string
	var lastReturnError error

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		if ticker == nil {
			ticker = time.NewTicker(duration)

			go func() {
				for {
					select {
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					case <-ctx.Done():
						ticker.Stop()
						break
					}
				}
			}()
		}

		if tokens > 0 {
			tokens--
			lastReturnString, lastReturnError = effector(ctx)
		}

		return lastReturnString, lastReturnError
	}
}
