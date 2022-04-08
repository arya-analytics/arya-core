package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"io"
)

func StreamDone(err error) (error, bool) {
	if err == io.EOF {
		return nil, true
	}
	if err != nil {
		return err, true
	}
	return nil, false
}

func StreamRange[T any](ctx context.Context, c chan T, f func(T) error) (err error) {
	route.RangeContext[T](ctx, c, func(v T) {
		if sErr, done := StreamDone(f(v)); done {
			err = sErr
			return
		}
	})
	return err
}

func StreamFor[T any](ctx context.Context, recv func() (T, error), action func(T) error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			v, rErr := recv()
			if err, done := StreamDone(rErr); done {
				return err
			}
			if err := action(v); err != nil {
				return err
			}
		}
	}
}
