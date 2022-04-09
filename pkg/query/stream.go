package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"io"
)

func StreamEOF(err error) (error, bool) {
	if err == io.EOF {
		return nil, true
	}
	return err, err != nil
}

func StreamDone(ctx context.Context, err error) (error, bool) {
	sErr, ok := StreamEOF(err)
	if ok || route.CtxDone(ctx) {
		return sErr, true
	}
	return sErr, ok
}

func StreamRange[T any](ctx context.Context, c chan T, f func(T) error) (err error) {
	route.RangeContext[T](ctx, c, func(v T) {
		if sErr, done := StreamEOF(f(v)); done {
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
			if err, done := StreamEOF(rErr); done {
				return err
			}
			if err := action(v); err != nil {
				return err
			}
		}
	}
}
