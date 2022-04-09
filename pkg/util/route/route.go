package route

import (
	"context"
)

func CtxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func RangeContext[T any](ctx context.Context, stream chan T, f func(T)) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-stream:
			if !ok {
				return
			}
			f(v)
		}
	}
}
