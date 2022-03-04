package chanchunk

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/google/uuid"
)

type QueryStreamCreateArgs struct {
	startTS telem.TimeStamp
	data    *telem.ChunkData
}

type QueryStreamCreate struct {
	obs        Observe
	exec       query.Execute
	rngSvc     *rng.Service
	configPK   uuid.UUID
	_config    *models.ChannelConfig
	_prevChunk *telem.Chunk
	prevCCPK   uuid.UUID
	errPipe    chan error
	catch      *errutil.CatchContext
	donePipe   chan bool
	stream     chan QueryStreamCreateArgs
	ctx        context.Context
}

func newStreamCreate(qExec query.Execute, obs Observe, rngSvc *rng.Service) *QueryStreamCreate {
	return &QueryStreamCreate{
		obs:      obs,
		exec:     qExec,
		rngSvc:   rngSvc,
		_config:  &models.ChannelConfig{},
		errPipe:  make(chan error, 10),
		stream:   make(chan QueryStreamCreateArgs),
		donePipe: make(chan bool),
	}
}

func (qsc *QueryStreamCreate) Start(ctx context.Context, pk uuid.UUID) error {
	qsc.ctx = ctx
	qsc.configPK = pk
	qsc.catch = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(qsc.errPipe)))
	if err := qsc.validateStart(); err != nil {
		return err
	}
	go qsc.listen()
	return nil
}

func (qsc *QueryStreamCreate) config() *models.ChannelConfig {
	if model.NewPK(qsc._config.ID).IsZero() {
		qsc.catch.Exec(query.NewRetrieve().BindExec(qsc.exec).Model(qsc._config).WherePK(qsc.configPK).Exec)
	}
	return qsc._config
}

func (qsc *QueryStreamCreate) updateConfigState(state models.ChannelState) error {
	qsc.obs.Add(ObservedChannelConfig{State: state, PK: qsc.configPK})
	qsc._config.State = state
	return query.NewUpdate().BindExec(qsc.exec).Model(qsc._config).WherePK(qsc.configPK).Fields("State").Exec(qsc.ctx)
}

func (qsc *QueryStreamCreate) validateStart() error {
	return validateStart().Exec(ValidateStartContext{cfg: qsc.config(), obs: qsc.obs}).Error()
}

func (qsc *QueryStreamCreate) prevChunk() *telem.Chunk {
	if qsc._prevChunk == nil {
		qsc.catch.Exec(func(ctx context.Context) error {
			ccr := &models.ChannelChunkReplica{}
			err := query.NewRetrieve().
				BindExec(qsc.exec).
				Model(ccr).
				Relation("ChannelChunk", "ID", "StartTS", "Size").
				WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": qsc.config().ID}).Exec(qsc.ctx)
			sErr, ok := err.(query.Error)
			if !ok || sErr.Type != query.ErrorTypeItemNotFound {
				return err
			}
			// If we don't find the item, this isn't an exceptional case, it just means the channel doesn't have any
			// data, so we can just return nil early.
			if sErr.Type == query.ErrorTypeItemNotFound {
				return nil
			}
			qsc._prevChunk = telem.NewChunk(ccr.ChannelChunk.StartTS, qsc.config().DataType, qsc.config().DataRate, ccr.Telem)
			return nil
		})
	}
	return qsc._prevChunk
}

func (qsc *QueryStreamCreate) setPrevChunk(chunk *telem.Chunk) {
	qsc._prevChunk = chunk
}

func (qsc *QueryStreamCreate) Send(startTS telem.TimeStamp, data *telem.ChunkData) {
	qsc.stream <- QueryStreamCreateArgs{startTS: startTS, data: data}
}

func (qsc *QueryStreamCreate) Close() {
	close(qsc.stream)
	<-qsc.donePipe
	close(qsc.errPipe)
}

func (qsc *QueryStreamCreate) Errors() chan error {
	return qsc.errPipe
}

func (qsc *QueryStreamCreate) listen() {
	qsc.catch.CatchSimple.Exec(func() error { return qsc.updateConfigState(models.ChannelStateActive) })
	for args := range qsc.stream {
		alloc := qsc.rngSvc.NewAllocate()

		nextChunk := telem.NewChunk(args.startTS, qsc.config().DataType, qsc.config().DataRate, args.data)
		qsc.catch.CatchSimple.Exec(func() error { return qsc.validateResolveNextChunk(nextChunk) })

		cc := &models.ChannelChunk{
			ID:              uuid.New(),
			ChannelConfigID: qsc.config().ID,
			StartTS:         nextChunk.Start(),
			Size:            nextChunk.Size(),
		}

		// CLARIFICATION: This means we tried to write a duplicate or consumed chunk.
		if cc.Size == 0 {
			continue
		}

		ccr := &models.ChannelChunkReplica{ID: uuid.New(), ChannelChunkID: cc.ID, Telem: args.data}

		qsc.catch.Exec(alloc.Chunk(qsc.config().NodeID, cc).Exec)
		qsc.catch.Exec(alloc.ChunkReplica(ccr).Exec)

		qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(cc).Exec)
		qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(ccr).Exec)

		qsc.setPrevChunk(nextChunk)
		qsc.catch.Reset()
	}
	qsc.catch.CatchSimple.Exec(func() error { return qsc.updateConfigState(models.ChannelStateInactive) })
	qsc.donePipe <- true
}

func (qsc *QueryStreamCreate) validateResolveNextChunk(nextChunk *telem.Chunk) error {
	vCtx := NextChunkValidateContext{prevChunk: qsc.prevChunk(), nextChunk: nextChunk}
	c := errutil.NewCatchSimple(errutil.WithAggregation())
	for _, err := range validateNextChunk().Exec(vCtx).Errors() {
		c.Exec(func() error {
			return qsc.resolveNextChunkError(err, vCtx)
		})
	}
	return c.Error()
}

func (qsc *QueryStreamCreate) resolveNextChunkError(err error, vCtx NextChunkValidateContext) error {
	return resolveNextChunk().Exec(err, NextChunkResolveContext{config: qsc.config(), NextChunkValidateContext: vCtx}).Error()
}

// |||| VALIDATE + RESOLVE ||||

// || START ||

type ValidateStartContext struct {
	obs Observe
	cfg *models.ChannelConfig
}

func validateStart() *validate.Validate[ValidateStartContext] {
	actions := []func(sCtx ValidateStartContext) error{validateConfigState}
	return validate.New(actions)
}

func validateConfigState(sCtx ValidateStartContext) error {
	oc, _ := sCtx.obs.Retrieve(sCtx.cfg.ID)
	if sCtx.cfg.State == models.ChannelStateActive || oc.State == models.ChannelStateActive {
		return errors.New("open a second stream to an active channel")
	}
	return nil
}

// || NEXT CHUNK ||

type NextChunkValidateContext struct {
	ctx       context.Context
	prevChunk *telem.Chunk
	nextChunk *telem.Chunk
}

func validateNextChunk() *validate.Validate[NextChunkValidateContext] {
	actions := []func(vCtx NextChunkValidateContext) error{validateTiming}
	return validate.New[NextChunkValidateContext](actions, errutil.WithAggregation())
}

type NextChunkResolveContext struct {
	config *models.ChannelConfig
	NextChunkValidateContext
}

func resolveNextChunk() *validate.Resolve[NextChunkResolveContext] {
	actions := []func(sErr error, rCtx NextChunkResolveContext) (bool, error){resolveTiming}
	return validate.NewResolve[NextChunkResolveContext](actions, errutil.WithAggregation())
}
