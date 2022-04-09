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

type StreamRetrieve struct {
	streamq.TSRetrieve
	qExec   query.Execute
	_config *models.ChannelConfig
	tRng    telem.TimeRange
	catch   *errutil.CatchContext
}

func newStreamRetrieve(qExec query.Execute) *StreamRetrieve {
	sr := &StreamRetrieve{qExec: qExec, _config: &models.ChannelConfig{}}
	sr.Base.Init(sr)
	sr.BindExec(sr.exec)
	return sr
}

func (sr *StreamRetrieve) WhereConfigPK(configPK uuid.UUID) *StreamRetrieve {
	newConfigPKOpt(sr.Pack(), configPK)
	return sr
}

func (sr *StreamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	sr.catch = errutil.NewCatchContext(context.Background())
	var (
		ccr        []*models.ChannelChunkReplica
		c          = *query.ConcreteModel[*chan *telem.Chunk](p)
		streamQ, _ = streamq.RetrieveStreamOpt(p, query.RequireOpt())
		tr, _      = streamq.RetrieveTimeRangeOpt(p, query.RequireOpt())
	)
	sr.catch.Exec(retrieveCCRQuery(sr.qExec, sr.config().ID, tr, &ccr).Exec)
	if sr.catch.Error() != nil {
		return sr.catch.Error()
	}
	streamQ.Segment(func() {
		defer func() {
			close(c)
			close(streamQ.Errors)
			streamQ.Complete()
		}()
		for _, r := range ccr {
			if route.CtxDone(ctx) {
				return
			}
			c <- telem.NewChunk(r.ChannelChunk.StartTS, sr.config().DataType, sr.config().DataRate, r.Telem)
		}
	})
	return nil
}

func (sr *StreamRetrieve) config() *models.ChannelConfig {
	configPK, _ := retrieveConfigPKOpt(sr.Pack(), query.RequireOpt())
	sr.catch.Exec(query.
		NewRetrieve().
		BindExec(sr.qExec).
		Model(sr._config).
		WherePK(configPK).
		WithMemo(query.NewMemo(sr._config)).
		Exec,
	)
	return sr._config
}

/// |||| QUERY UTILS ||||

func retrieveCCRQuery(
	qExec query.Execute,
	configPK uuid.UUID,
	tr telem.TimeRange,
	ccr *[]*models.ChannelChunkReplica,
) *query.Retrieve {
	return query.NewRetrieve().
		BindExec(qExec).
		Model(ccr).
		Relation("ChannelChunk", "StartTS").
		WhereFields(query.WhereFields{
			"ChannelChunk.ChannelConfigID": configPK,
			"ChannelChunk.StartTS":         query.InRange(tr.Start(), tr.End()),
		})
}

// |||| OPTS ||||

const configPKOptKey query.OptKey = "configPK"

func newConfigPKOpt(p *query.Pack, configPK uuid.UUID) {
	p.SetOpt(configPKOptKey, configPK)
}

func retrieveConfigPKOpt(p *query.Pack, opts ...query.OptRetrieveOpt) (uuid.UUID, bool) {
	o, ok := p.RetrieveOpt(configPKOptKey, opts...)
	if !ok {
		return uuid.Nil, false
	}
	return o.(uuid.UUID), true
}
