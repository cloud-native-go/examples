package ch04

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestDebounceFirstDataRace tests for data races.
func TestDebounceFirstDataRace(t *testing.T) {
	ctx := context.Background()

	circuit := failAfter(1)
	debounce := DebounceFirst(circuit, time.Second)

	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}
