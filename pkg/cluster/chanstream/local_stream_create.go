package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type localStreamCreate struct {
	qe query.Execute
}

func newLocalStreamCreate(qe query.Execute) *localStreamCreate {
	return &localStreamCreate{qe: qe}
}

func (lsc *localStreamCreate) exec(ctx context.Context, p *query.Pack) error {
	sampleStream := *query.ConcreteModel[*chan *models.ChannelSample](p)
	streamQ, _ := streamq.RetrieveStreamOpt(p, query.RequireOpt())
	streamQ.Segment(func() {
		for s := range sampleStream {
			if route.CtxDone(ctx) {
				return
			}
			if err := streamq.NewTSCreate().Model(s).BindExec(lsc.qe).Exec(ctx); err != nil {
				streamQ.Errors <- err
			}
		}
	}, streamq.WithSegmentName("cluster.telemstream.localStreamCreate"))
	return nil
}
