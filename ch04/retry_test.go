package ch04

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

var count int

func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

func TestRetry(t *testing.T) {
	ctx := context.Background()
	r := Retry(EmulateTransientError, 5, 2*time.Second)
	res, err := r(ctx)

	fmt.Println(res, err)
}
