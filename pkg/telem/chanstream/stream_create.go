package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type streamCreate struct {
	qExec query.Execute
}

const errorPipeCapacity = 10

func newStreamCreate(qExec query.Execute) *streamCreate {
	return &streamCreate{qExec: qExec}
}

func (sc *streamCreate) exec(ctx context.Context, p *query.Pack) error {
	return sc.qExec(ctx, p)
}
