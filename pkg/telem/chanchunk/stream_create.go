package chanchunk

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/google/uuid"
)

type StreamCreate struct {
	streamq.TSCreate
	obs          observe
	qExec        query.Execute
	rngSvc       *rng.Service
	_config      *models.ChannelConfig
	_prevChunk   *telem.Chunk
	prevCCPK     uuid.UUID
	catch        *errutil.CatchContext
	streamQ      *streamq.Stream
	valStream    chan StreamCreateArgs
	cancelStream context.CancelFunc
}

type StreamCreateArgs struct {
	Start telem.TimeStamp
	Data  *telem.ChunkData
}

func newStreamCreate(qExec query.Execute, obs observe, rngSvc *rng.Service) *StreamCreate {
	sc := &StreamCreate{qExec: qExec, obs: obs, rngSvc: rngSvc, _config: &models.ChannelConfig{}}
	sc.Base.Init(sc)
	sc.BindExec(sc.exec)
	return sc
}

func (sc *StreamCreate) WhereConfigPK(configPK uuid.UUID) *StreamCreate {
	newConfigPKOpt(sc.Pack(), configPK)
	return sc
}

func (sc *StreamCreate) exec(ctx context.Context, p *query.Pack) error {
	sc.valStream = *query.ConcreteModel[*chan StreamCreateArgs](p)
	sc.streamQ, _ = streamq.RetrieveStreamOpt(p, query.RequireOpt())
	sc.catch = errutil.NewCatchContext(context.Background(), errutil.WithHooks(errutil.NewPipeHook(sc.streamQ.Errors)))
	if err := sc.validateStart(); err != nil {
		return err
	}
	sc.listen(ctx)
	return nil
}

// |||| PROCESS ||||

func (sc *StreamCreate) listen(ctx context.Context) {
	sc.streamQ.Segment(func() {
		sc.updateConfigStatus(models.ChannelStatusActive)
		defer sc.streamQ.Complete()
		defer sc.updateConfigStatus(models.ChannelStatusInactive)
		route.RangeContext(ctx, sc.valStream, sc.processNextChunk)
	}, streamq.WithSegmentName("telem.chanchunk.StreamCreate"))
}

func (sc *StreamCreate) processNextChunk(args StreamCreateArgs) {
	nc := telem.NewChunk(args.Start, sc.config().DataType, sc.config().DataRate, args.Data)
	sc.validateResolveNextChunk(nc)

	cc := &models.ChannelChunk{
		ID:              uuid.New(),
		ChannelConfigID: sc.config().ID,
		StartTS:         nc.Start(),
		Size:            nc.Size(),
	}

	ccr := &models.ChannelChunkReplica{
		ID:             uuid.New(),
		ChannelChunkID: cc.ID,
		Telem:          nc.ChunkData,
	}

	// CLARIFICATION: This means we tried to write a duplicate or consumed chunk.
	if cc.Size == 0 {
		return
	}

	a := sc.rngSvc.NewAllocate()
	sc.catch.Exec(a.Chunk(sc.config().NodeID, cc).Exec)
	sc.catch.Exec(a.ChunkReplica(ccr).Exec)

	sc.catch.Exec(query.NewCreate().BindExec(sc.qExec).Model(cc).Exec)
	sc.catch.Exec(query.NewCreate().BindExec(sc.qExec).Model(ccr).Exec)

	sc.setPrevChunk(nc)
	sc.catch.Reset()
}

func (sc *StreamCreate) config() *models.ChannelConfig {
	configPK, _ := retrieveConfigPKOpt(sc.Pack(), query.RequireOpt())
	sc.catch.Exec(query.
		NewRetrieve().
		BindExec(sc.qExec).
		Model(sc._config).
		WherePK(configPK).
		WithMemo(query.NewMemo(sc._config)).
		Exec,
	)
	return sc._config
}

func (sc *StreamCreate) updateConfigStatus(status models.ChannelState) {
	sc.obs.Add(observedChannelConfig{State: status, PK: sc.config().ID})
	sc.config().State = status
	sc.catch.CatchSimple.Exec(func() error {
		return query.NewUpdate().
			BindExec(sc.qExec).
			Model(sc.config()).
			WherePK(sc.config().ID).
			Fields("State").Exec(context.Background())
	})
}

func (sc *StreamCreate) prevChunk() *telem.Chunk {
	if sc._prevChunk == nil {
		sc.catch.Exec(func(ctx context.Context) error {
			ccr := &models.ChannelChunkReplica{}
			err := query.NewRetrieve().
				BindExec(sc.qExec).
				Model(ccr).
				Relation("ChannelChunk", "ID", "StartTS", "Size").
				WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": sc.config().ID}).Exec(ctx)
			sErr, ok := err.(query.Error)
			if !ok || sErr.Type != query.ErrorTypeItemNotFound {
				return err
			}
			// If we don't find the item, this isn't an exceptional case, it just means the channel doesn't have any
			// Data, so we can just return nil early.
			if sErr.Type == query.ErrorTypeItemNotFound {
				return nil
			}
			sc._prevChunk = telem.NewChunk(ccr.ChannelChunk.StartTS, sc.config().DataType, sc.config().DataRate, ccr.Telem)
			return nil
		})
	}
	return sc._prevChunk
}

func (sc *StreamCreate) setPrevChunk(chunk *telem.Chunk) {
	sc._prevChunk = chunk
}

// |||| VALIDATE + RESOLVE ||||

func (sc *StreamCreate) validateStart() error {
	return validateStart().Exec(validateStartContext{cfg: sc.config(), obs: sc.obs}).Error()
}
func (sc *StreamCreate) validateResolveNextChunk(nextChunk *telem.Chunk) {
	nc := nextChunkContext{cfg: sc.config(), prev: sc.prevChunk(), next: nextChunk}
	sc.catch.CatchSimple.Exec(func() error {
		for _, vErr := range validateNextChunk().Exec(nc).Errors() {
			if rErr := sc.resolveNextChunkError(vErr, nc); rErr != nil {
				return rErr
			}
		}
		return nil
	})
}

func (sc *StreamCreate) resolveNextChunkError(err error, nCtx nextChunkContext) error {
	return resolveNextChunk().Exec(err, nCtx).Error()
}

// || START ||

type validateStartContext struct {
	obs observe
	cfg *models.ChannelConfig
}

func validateStart() *validate.Validate[validateStartContext] {
	actions := []func(sCtx validateStartContext) error{validateConfigState}
	return validate.New(actions)
}

func validateConfigState(sCtx validateStartContext) error {
	oc, _ := sCtx.obs.Retrieve(sCtx.cfg.ID)
	if sCtx.cfg.State == models.ChannelStatusActive || oc.State == models.ChannelStatusActive {
		return query.NewSimpleError(
			query.ErrorTypeInvalidArgs,
			errors.New("cannot open a second stream on an active channel"),
		)

	}
	return nil
}

// || NEXT CHUNK ||

type nextChunkContext struct {
	ctx  context.Context
	cfg  *models.ChannelConfig
	prev *telem.Chunk
	next *telem.Chunk
}

func validateNextChunk() *validate.Validate[nextChunkContext] {
	a := []func(vCtx nextChunkContext) error{validateTiming}
	return validate.New[nextChunkContext](a, errutil.WithAggregation())
}

func resolveNextChunk() *validate.Resolve[nextChunkContext] {
	a := []func(sErr error, rCtx nextChunkContext) (bool, error){resolveTiming}
	return validate.NewResolve[nextChunkContext](a, errutil.WithAggregation())
}
