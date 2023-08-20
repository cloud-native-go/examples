package ch04

import (
	"context"
	"sync"
	"time"
)

func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	var m sync.Mutex
	var timer *time.Timer
	var cctx context.Context
	var cancel context.CancelFunc

	return func(ctx context.Context) (string, error) {
		m.Lock()

		if timer != nil {
			timer.Stop()
			cancel()
		}

		cctx, cancel = context.WithCancel(ctx)
		ch := make(chan struct {
			result string
			err    error
		}, 1)

		timer = time.AfterFunc(d, func() {
			r, e := circuit(cctx)
			ch <- struct {
				result string
				err    error
			}{r, e}
		})

		m.Unlock()

		select {
		case res := <-ch:
			return res.result, res.err
		case <-cctx.Done():
			return "", cctx.Err()
		}
	}
}
