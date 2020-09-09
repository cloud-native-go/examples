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
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Number of concurrently requesting instances
const instanceCount = 1000

// How long to run each trial
const trialDuration = 4 * time.Minute

// How much time is allocated to each bucket
const bucketWidth = time.Second

// Slice to track request counts
var requestBuckets []int

// The index of the current bucket
var currentBucketIndex int

// An "events" channel. It is used by sendRequest.
var requestEvents chan bool = make(chan bool)

// The backoff function to use
var backoffFunction func() string = withExponentialBackoffAndJitter

// Each instance will randomly delay up to this duration
var maxStartupDelay = bucketWidth

// Just used to track the time the program started, for output purposes.
var startTime = time.Now()

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	bucketCount := int(trialDuration / bucketWidth)
	requestBuckets = make([]int, bucketCount)
	log("Bucket count: %d\n", bucketCount)

	go catchEvents()

	log("Starting %d backoff processes\n", instanceCount)
	for i := 0; i < instanceCount; i++ {
		go func() {
			delay := time.Duration(rand.Int63n(int64(maxStartupDelay)))
			time.Sleep(delay)
			backoffFunction()
		}()
	}
	log("%d backoff processes started\n", instanceCount)

	for currentBucketIndex = 0; currentBucketIndex < bucketCount; currentBucketIndex++ {
		time.Sleep(bucketWidth)

		i := currentBucketIndex
		if i >= bucketCount {
			i = bucketCount - 1
		}

		log("Bucket %d: %d data points\n", currentBucketIndex, requestBuckets[i])
	}

	sum := 0
	for i := 0; i < bucketCount; i++ {
		sum += requestBuckets[i]
		fmt.Println(requestBuckets[i])
	}

	log("Avg: %d\n", sum/bucketCount)
}

// Using withNoBackoff as the backoff function sends retries as quickly as
// possible, with no backoff of any kind.
func withNoBackoff() string {
	res, err := sendRequest()
	for err != nil {
		res, err = sendRequest()
	}

	return res
}

// Using withDelayedBackoff as the backoff function sends retries after a
// two-second delay.
func withDelayedBackoff() string {
	res, err := sendRequest()
	for err != nil {
		time.Sleep(2000 * time.Millisecond)
		res, err = sendRequest()
	}

	return res
}

// Using withExponentialBackoff as the backoff function sends retries
// initially with a 1-second delay, but doubling after each attempt to
// a maximum delay of 1-minute.
func withExponentialBackoff() string {
	res, err := sendRequest()
	base, cap := time.Second, time.Minute

	for backoff := base; err != nil; backoff <<= 1 {
		if backoff > cap {
			backoff = cap
		}
		time.Sleep(backoff)
		res, err = sendRequest()
	}

	return res
}

// Using withExponentialBackoff as the backoff function sends retries
// initially with a 1-second delay, but doubling after each attempt to
// a maximum delay of 1-minute. Each backoff time gets an additional
func withExponentialBackoffAndJitter() string {
	res, err := sendRequest()
	base, cap := time.Second, time.Minute

	for backoff := base; err != nil; backoff <<= 1 {
		if backoff > cap {
			backoff = cap
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := base + time.Duration(jitter)
		time.Sleep(sleep)
		res, err = sendRequest()
	}

	return res
}

// sendRequest simulates sending a request. It always returns an
// error after a short delay.
func sendRequest() (string, error) {
	delay := time.Duration(rand.Int63n(100)+rand.Int63n(100)) * time.Millisecond

	time.Sleep(delay / 2)
	requestEvents <- true
	time.Sleep(delay / 2)

	return "", errors.New("")
}

// catchEvents listens on the requestEvents channel and increments the
// appropriate bucket.
func catchEvents() {
	for range requestEvents {
		requestBuckets[currentBucketIndex]++
	}
}

// log emits timestamped log output.
func log(f string, i ...interface{}) {
	since := time.Now().Sub(startTime)
	t := time.Time{}.Add(since)
	tf := t.Format("15:04:05")

	fmt.Printf(tf+" "+f, i...)
}
