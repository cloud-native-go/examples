package ch04

import "context"

type TimeoutFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

func Timeout(f TimeoutFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		ch := make(chan struct {
			result string
			err    error
		}, 1)

		go func() {
			res, err := f(arg)
			ch <- struct {
				result string
				err    error
			}{res, err}
		}()

		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
