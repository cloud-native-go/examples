/*
 * Copyright 2020 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Effector is the function that you want to subject to throttling.
type Effector func(context.Context) (string, error)

// Throttled wraps an Effector. It accepts the same parameters, plus a
// "key" string that represents a caller identity. It returns the same,
// plus a bool that's true if the call is not throttled.
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
			return true, str, err
		}

		// Calculate how many tokens we now have based on the time
		// passed since the previous request.
		refillEventsSince := uint(time.Since(r.time) / d)
		tokensAddedSince := refill * refillEventsSince
		currentTokens := r.tokens + tokensAddedSince

		// We don't have enough tokens. Return false.
		if currentTokens < 1 {
			return false, "", nil
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

		return true, str, err
	}
}

func getHostname(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return os.Hostname()
}

func main() {
	throttled := Throttle(getHostname, 2, 2, 500*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for range ticker.C {
		if ok, msg, err := throttled(ctx, "foo"); err != nil {
			fmt.Println(err)
			break
		} else if ok {
			fmt.Println("OK! Got", msg)
		} else {
			fmt.Println("Throttled :(")
		}
	}
}
