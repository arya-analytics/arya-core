package chanchunk

import (
	"context"
	"errors"
	"fmt"
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

type streamCreate struct {
	obs        observe
	qExec      query.Execute
	rngSvc     *rng.Service
	configPK   uuid.UUID
	_config    *models.ChannelConfig
	_prevChunk *telem.Chunk
	prevCCPK   uuid.UUID
	catch      *errutil.CatchContext
	streamQ    *streamq.Stream
	valStream  chan StreamCreateArgs
	ctx        context.Context
}

type StreamCreateArgs struct {
	Start telem.TimeStamp
	Data  *telem.ChunkData
}

func newStreamCreate(qExec query.Execute, obs observe, rngSvc *rng.Service) *streamCreate {
	return &streamCreate{qExec: qExec, obs: obs, rngSvc: rngSvc}
}

func (sc *streamCreate) exec(ctx context.Context, p *query.Pack) error {
	sc.ctx = ctx
	pkc, _ := query.PKOpt(p, query.PanicIfOptNotPresent())
	if len(pkc) != 1 {
		panic(fmt.Sprintf("stream_create: expected 1 pk, got %v", len(pkc)))
	}
	sc.configPK = pkc[0].Raw().(uuid.UUID)
	streamQ, _ := streamq.StreamOpt(p, query.PanicIfOptNotPresent())
	sc.catch = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(streamQ.Errors)))
	if err := sc.validateStart(); err != nil {
		return err
	}
	sc.listen()
	return nil
}

// |||| PROCESS ||||

func (sc *streamCreate) listen() {
	sc.streamQ.Segment(func() {
		sc.updateConfigStatus(models.ChannelStatusActive)
		defer sc.updateConfigStatus(models.ChannelStatusInactive)
		for args := range sc.valStream {
			if route.CtxDone(sc.ctx) {
				return
			}
			sc.processNextChunk(args.Start, args.Data)
		}
	}, streamq.WithSegmentName("telem.chanchunk.streamCreate"))
}

func (sc *streamCreate) processNextChunk(startTS telem.TimeStamp, data *telem.ChunkData) {
	nc := telem.NewChunk(startTS, sc.config().DataType, sc.config().DataRate, data)
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
		Telem:          data,
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

func (sc *streamCreate) config() *models.ChannelConfig {
	sc.catch.Exec(query.
		NewRetrieve().
		BindExec(sc.qExec).
		Model(sc._config).
		WherePK(sc.configPK).
		WithMemo(query.NewMemo(sc._config)).
		Exec,
	)
	return sc._config
}

func (sc *streamCreate) updateConfigStatus(status models.ChannelStatus) {
	sc.obs.Add(observedChannelConfig{Status: status, PK: sc.configPK})
	sc.config().Status = status
	sc.catch.CatchSimple.Exec(func() error {
		return query.NewUpdate().
			BindExec(sc.qExec).
			Model(sc.config()).
			WherePK(sc.configPK).
			Fields("Status").Exec(context.Background())
	})
}

func (sc *streamCreate) prevChunk() *telem.Chunk {
	if sc._prevChunk == nil {
		sc.catch.Exec(func(ctx context.Context) error {
			ccr := &models.ChannelChunkReplica{}
			err := query.NewRetrieve().
				BindExec(sc.qExec).
				Model(ccr).
				Relation("ChannelChunk", "ID", "StartTS", "Size").
				WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": sc.config().ID}).Exec(sc.ctx)
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

func (sc *streamCreate) setPrevChunk(chunk *telem.Chunk) {
	sc._prevChunk = chunk
}

// |||| VALIDATE + RESOLVE ||||

func (sc *streamCreate) validateStart() error {
	return validateStart().Exec(validateStartContext{cfg: sc.config(), obs: sc.obs}).Error()
}
func (sc *streamCreate) validateResolveNextChunk(nextChunk *telem.Chunk) {
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

func (sc *streamCreate) resolveNextChunkError(err error, nCtx nextChunkContext) error {
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
	if sCtx.cfg.Status == models.ChannelStatusActive || oc.Status == models.ChannelStatusActive {
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
