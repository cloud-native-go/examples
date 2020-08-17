package main

import (
	"context"
	"fmt"
	"time"
)

// Effector is the function that you want to subject to throttling.
type Effector func(context.Context) (string, error)

// Throttled wraps an Effector. It accepts the same parameters, plus a
// "key" string that represents a caller identity. It returns the same,
// plus a bool that's true if the call is throttled.
type Throttled func(context.Context, string) (bool, string, error)

type record struct {
	tokens uint
	time   time.Time
}

// Throttle accepts an Effector function, and returns a Throttled
// function with a per-key token bucket with a capacity of max
// that refills at a rate of refill tokens every d.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Throttled {
	bucket := map[string]*record{}

	return func(ctx context.Context, key string) (bool, string, error) {
		r := bucket[key]

		// This is a new entry! It passes. Naïvely assumes that
		// Capacity >= TokensConsumedPerRequest.
		if r == nil {
			tokens := max - 1
			bucket[key] = &record{tokens: tokens, time: time.Now()}

			str, err := e(ctx)
			return false, str, err
		}

		// Calculate how many tokens we now have based on the time
		// passed since the previous request.
		refillEventsSince := uint(time.Since(r.time) / d)
		tokensAddedSince := refill * refillEventsSince
		currentTokens := r.tokens + tokensAddedSince

		// We don't have enough tokens. Return false.
		if currentTokens < 1 {
			return true, "", nil
		}

		// If we've refilled our bucket, we can restart the clock.
		// Otherwise, we figure out when the most recent tokens were added.
		if currentTokens > max {
			r.time = time.Now()
			r.tokens = max - 1
		} else {
			deltaTokens := currentTokens - r.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d

			r.time = r.time.Add(deltaTime)
			r.tokens = currentTokens - 1
		}

		str, err := e(ctx)

		return false, str, err
	}
}

func doSomething(ctx context.Context) (string, error) {
	return "something", nil
}

func main() {
	throttled := Throttle(doSomething, 2, 2, 500*time.Millisecond)

	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for range ticker.C {
		throttled, msg, _ := throttled(context.Background(), "foo")

		if throttled {
			fmt.Println("Throttled :(")
		} else {
			fmt.Println("OK! Got", msg)
		}
	}
}
