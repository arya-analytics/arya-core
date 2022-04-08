package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

type streamRetrieve struct {
	qExec    query.Execute
	configPK uuid.UUID
	_config  *models.ChannelConfig
	tRng     telem.TimeRange
	catch    *errutil.CatchContext
}

func newStreamRetrieve(qExec query.Execute) *streamRetrieve {
	return &streamRetrieve{qExec: qExec, _config: &models.ChannelConfig{}}
}

func (sr *streamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	pkc, _ := query.RetrievePKOpt(p, query.RequireOpt())
	sr.configPK = pkc[0].Raw().(uuid.UUID)
	sr.catch = errutil.NewCatchContext(context.Background())
	var (
		replicas   []*models.ChannelChunkReplica
		c          = *query.ConcreteModel[*chan *telem.Chunk](p)
		streamQ, _ = streamq.RetrieveStreamOpt(p, query.RequireOpt())
		tRng, _    = streamq.RetrieveTimeRangeOpt(p, query.RequireOpt())
	)
	sr.catch.Exec(query.NewRetrieve().
		BindExec(sr.qExec).
		Model(&replicas).
		Relation("ChannelChunk", "StartTS").
		WhereFields(query.WhereFields{
			"ChannelChunk.StartTS":         query.InRange(tRng.Start(), tRng.End()),
			"ChannelChunk.ChannelConfigID": sr.config().ID,
		}).
		Exec)
	if sr.catch.Error() != nil {
		return sr.catch.Error()
	}
	streamQ.Segment(func() {
		defer close(c)
		for _, r := range replicas {
			if route.CtxDone(ctx) {
				return
			}
			c <- telem.NewChunk(r.ChannelChunk.StartTS, sr.config().DataType, sr.config().DataRate, r.Telem)
		}
	})
	return nil
}

func (sr *streamRetrieve) config() *models.ChannelConfig {
	sr.catch.Exec(query.
		NewRetrieve().
		BindExec(sr.qExec).
		Model(sr._config).
		WherePK(sr.configPK).
		WithMemo(query.NewMemo(sr._config)).
		Exec,
	)
	return sr._config
}
