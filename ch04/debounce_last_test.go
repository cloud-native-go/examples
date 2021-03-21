package ch04

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestDebounceLastDataRace tests for data races.
func TestDebounceLastDataRace(t *testing.T) {
	ctx := context.Background()
	debounce := DebounceLast(counter(), time.Second)
	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()

	t.Log("Waiting 2 seconds")

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()
}
