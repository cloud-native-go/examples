/*
 * Copyright 2024 Matthew A. Titmus
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

package ch04

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestFuture just runs slow functions, and makes sure that it returns the
// expected result after the expected amount of time.
func TestFuture(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)
}

// TestFutureGetTwice tests that subsequent calls to future.Result()
// immediately return the initial return values.
func TestFutureGetTwice(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)

	// Get result again. Should happen straightaway.

	start = time.Now()

	res, err = future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 0)
}

// TestFutureConcurrent tests that the Future is thread-safe.
func TestFutureConcurrent(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			res, err := future.Result()
			if err != nil {
				t.Error(err)
				return
			}

			if !strings.HasPrefix(res, "I slept for") {
				t.Error("unexpected output:", res)
			}

			elapsedCheck(t, start, 2)
		}()
	}

	wg.Wait()
}

// TestFutureTimeout makes sure that the future will time out with an error
// if its context is canceled with a timeout
func TestFutureTimeout(t *testing.T) {
	start := time.Now()

	// Get a context decorated with a 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	// The cancel function returned by context.WithTimeout should be called,
	// not discarded, to avoid a context leak
	defer cancel()

	future := SlowFunction(ctx)

	// We should time out with a "context deadline exceeded" error
	res, err := future.Result()
	if err != nil {
		if !strings.Contains(err.Error(), "deadline") {
			t.Error("received unexpected error maybe: ", err)
		}
	}

	// Result should be empty
	if res != "" {
		t.Error("should have an empty result")
	}

	// Timeout should be after 1 second
	elapsedCheck(t, start, 1)
}

// TestFutureCancel
func TestFutureCancel(t *testing.T) {
	start := time.Now()

	// Get a context with an explicit cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Wait a second, and then cancel the future.
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	future := SlowFunction(ctx)

	// We should time out with a "context deadline exceeded" error
	res, err := future.Result()
	if err != nil {
		if !strings.Contains(err.Error(), "canceled") {
			t.Error("received unexpected error maybe: ", err)
		}
	}

	// Result should be empty
	if res != "" {
		t.Error("should have an empty result")
	}

	// Timeout should be after 1 second
	elapsedCheck(t, start, 1)
}

func elapsedCheck(t *testing.T, start time.Time, seconds int) {
	elapsed := int(time.Now().Sub(start).Seconds())

	if seconds != elapsed {
		t.Errorf("expected %d seconds; got %d\n", seconds, elapsed)
	}
}
