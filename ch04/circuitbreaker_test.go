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
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func counter() Circuit {
	m := sync.Mutex{}
	count := 0

	// Service function. Fails after 5 tries.
	return func(ctx context.Context) (string, error) {
		m.Lock()
		count++
		m.Unlock()

		return fmt.Sprintf("%d", count), nil
	}
}

// failAfter returns a function matching the Circuit type that returns an
// error after its been called more than threshold times.
func failAfter(threshold int) Circuit {
	count := 0

	// Service function. Fails after 5 tries.
	return func(ctx context.Context) (string, error) {
		count++

		if count > threshold {
			return "", errors.New("INTENTIONAL FAIL!")
		}

		return "Success", nil
	}
}

func waitAndContinue() Circuit {
	return func(ctx context.Context) (string, error) {
		time.Sleep(time.Second)

		if rand.Int()%2 == 0 {
			return "success", nil
		}

		return "Failed", fmt.Errorf("forced failure")
	}
}

// TestFailAfter5 tests that the failAfter5 test function acts as expected.
func TestCircuitBreakerFailAfter5(t *testing.T) {
	circuit := failAfter(5)
	ctx := context.Background()

	for count := 1; count <= 5; count++ {
		_, err := circuit(ctx)

		t.Logf("attempt %d: %v", count, err)

		switch {
		case count <= 5 && err != nil:
			t.Error("expected no error; got", err)
		case count > 5 && err == nil:
			t.Error("expected err; got none")
		}
	}
}

// TestBreaker tests that the Breaker function automatically closes and reopens.
func TestCircuitBreaker(t *testing.T) {
	// Service function. Fails after 5 tries.
	circuit := failAfter(5)

	// A circuit breaker that opens after one failed attempt.
	breaker := Breaker(circuit, 1)

	ctx := context.Background()

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReclose := false
	count := 0

	for range time.NewTicker(time.Second).C {
		_, err := breaker(ctx)

		if err != nil {
			// Does the circuit open?
			if strings.HasPrefix(err.Error(), "service unreachable") {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true

					t.Log("circuit has opened")
				}
			} else {
				// Does it close again?
				if circuitOpen {
					circuitOpen = false
					doesCircuitReclose = true

					t.Log("circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit closed and operational")
		}

		count++
		if count >= 10 {
			break
		}
	}

	if !doesCircuitOpen {
		t.Error("circuit didn't appear to open")
	}

	if !doesCircuitReclose {
		t.Error("circuit didn't appear to close after time")
	}
}

// TestCircuitBreakerDataRace tests for data races.
func TestCircuitBreakerDataRace(t *testing.T) {
	ctx := context.Background()

	circuit := waitAndContinue()
	breaker := Breaker(circuit, 1)

	wg := sync.WaitGroup{}

	for count := 1; count <= 20; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := breaker(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}
