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
	"fmt"
	"testing"
	"time"
)

// callsCountFunction returns a function that increments the int that
// callCounter points to whenever it is run.
func callsCountFunction(callCounter *int) Effector {
	return func(ctx context.Context) (string, error) {
		*callCounter++
		return fmt.Sprintf("call %d", *callCounter), nil
	}
}

// TestThrottleMax1 tests whether a max of 1 call per duration is respected.
func TestThrottleMax1(t *testing.T) {
	const max uint = 1

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}

	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

// TestThrottleMax10 tests whether a max of 10 calls per duration is respected.
func TestThrottleMax10(t *testing.T) {
	const max uint = 10

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}

	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

// TestThrottleCallFrequency5Seconds tests whether a Throttle with a max of 1
// and a duration of 1 second called every 250ms for 5 seconds will be called
// exactly 5 times.
func TestThrottleCallFrequency5Seconds(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 1, 1, time.Second)

	// make a call every 1/4 second for 5 seconds.
	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond).C

	for range ticker {
		tickCounts++

		s, e := throttle(ctx)
		if e != nil {
			t.Log("Error:", e)
		} else {
			t.Log("output:", s)
		}

		if tickCounts >= 20 {
			break
		}
	}

	if callsCounter != 5 {
		t.Error("expected 5; got", callsCounter)
	}
}

// TestThrottleVariableRefill
func TestThrottleVariableRefill(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 4, 2, 500*time.Millisecond)

	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond)
	timer := time.NewTimer(2 * time.Second)

time:
	for {
		select {
		case <-ticker.C:
			tickCounts++

			s, e := throttle(ctx)
			if e != nil {
				t.Log("Error:", e)
			} else {
				t.Log("output:", s)
			}
		case <-timer.C:
			break time
		}
	}

	if callsCounter != 8 {
		t.Error("expected 8; got", callsCounter)
	}
}

// TestThrottleContextTimeout tests whether a Throttle will return an error
// when its context is canceled.
func TestThrottleContextTimeout(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	throttle := Throttle(effector, 1, 1, time.Second)

	s, e := throttle(ctx)
	if e != nil {
		t.Error("unexpected error:", e)
	} else {
		t.Log("output:", s)
	}

	// Wait for timeout
	time.Sleep(300 * time.Millisecond)

	_, e = throttle(ctx)
	if e != nil {
		t.Log("got expected error:", e)
	} else {
		t.Error("didn't get expected error")
	}
}
