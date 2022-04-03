package query

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"io"
)

func StreamDone(ctx context.Context, err error) (error, bool) {
	if route.CtxDone(ctx) || err == io.EOF {
		return nil, true
	}
	if err != nil {
		return err, true
	}
	return nil, false
}
