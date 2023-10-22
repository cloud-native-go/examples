/*
 * Copyright 2023 Matthew A. Titmus
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
	"fmt"
	"testing"
	"time"
)

func TestChord(t *testing.T) {
	expected := [][]int{
		{2, 1, 1},
		{4, 2, 2},
		{6, 3, 3},
		{8, 4, 4},
		{10, 5, 5},
	}

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		for n := 1; n <= 5; n++ {
			ch1 <- n
			ch1 <- n * 2
			ch2 <- n
			ch3 <- n
			time.Sleep(time.Second)
		}

		// Chord ends when any channel closes
		close(ch1)
		close(ch2)
		close(ch3)
	}()

	res := Chord(ch1, ch2, ch3)
	idx := 0

	for ss := range res {
		fmt.Println(ss)

		for i, s := range ss {
			if expected[idx][i] != s {
				t.Fatalf("Expected: %v; Got: %v", expected[idx], ss)
			}
		}

		idx++
	}
}
