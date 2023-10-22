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
