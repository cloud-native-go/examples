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
	"testing"
	"time"
)

func TestTimeoutNo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	timeout := Timeout(Slow)
	_, err := timeout(ctx, "some input")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestTimeoutYes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	defer cancel()

	timeout := Timeout(Slow)
	_, err := timeout(ctx, "some input")
	if err == nil {
		t.Fatal("Didn't get expected timeout error")
	}
}

func Slow(s string) (string, error) {
	time.Sleep(time.Second)
	return "Got input: " + s, nil
}
