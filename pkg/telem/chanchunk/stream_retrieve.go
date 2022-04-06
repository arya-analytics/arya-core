package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

type StreamRetrieve struct {
	qa    query.Assemble
	cfgPK interface{}
	_cfg  *models.ChannelConfig
	tRng  telem.TimeRange
	catch *errutil.CatchContext
}

func newStreamRetrieve(qa query.Assemble) *StreamRetrieve {
	return &StreamRetrieve{qa: qa, _cfg: &models.ChannelConfig{}}
}

func (sr *StreamRetrieve) WhereConfigPK(cfgPK interface{}) *StreamRetrieve {
	sr.cfgPK = cfgPK
	return sr
}

func (sr *StreamRetrieve) WhereTimeRange(tRng telem.TimeRange) *StreamRetrieve {
	sr.tRng = tRng
	return sr
}

func (sr *StreamRetrieve) config() *models.ChannelConfig {
	if model.NewPK(sr._cfg.ID).IsZero() {
		sr.catch.Exec(sr.qa.NewRetrieve().Model(sr._cfg).WherePK(sr.cfgPK).Exec)
	}
	return sr._cfg
}

func (sr *StreamRetrieve) Exec(ctx context.Context) (chan *telem.Chunk, error) {
	sr.catch = errutil.NewCatchContext(ctx)
	var (
		stream   = make(chan *telem.Chunk)
		replicas []*models.ChannelChunkReplica
	)
	cfg := sr.config()
	sr.catch.Exec(sr.qa.NewRetrieve().
		Model(&replicas).
		Relation("ChannelChunk", "StartTS").
		WhereFields(query.WhereFields{
			"ChannelChunk.StartTS":          query.InRange(sr.tRng.Start(), sr.tRng.End()),
			"ChannelChunk.ChannelConfig.ID": cfg.ID,
		}).
		Exec)
	if sr.catch.Error() != nil {
		return stream, sr.catch.Error()
	}
	go func() {
		for _, r := range replicas {
			stream <- telem.NewChunk(
				r.ChannelChunk.StartTS,
				cfg.DataType,
				cfg.DataRate,
				r.Telem,
			)
		}
		close(stream)
	}()
	return stream, nil
}
