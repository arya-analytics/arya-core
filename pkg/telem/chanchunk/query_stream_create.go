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
	stream     chan queryStreamCreateArgs
	ctx        context.Context
}

type queryStreamCreateArgs struct {
	startTS telem.TimeStamp
	data    *telem.ChunkData
}

func newStreamCreate(qExec query.Execute, obs Observe, rngSvc *rng.Service) *QueryStreamCreate {
	return &QueryStreamCreate{
		obs:      obs,
		exec:     qExec,
		rngSvc:   rngSvc,
		_config:  &models.ChannelConfig{},
		errPipe:  make(chan error, 10),
		stream:   make(chan queryStreamCreateArgs),
		donePipe: make(chan bool),
	}
}

func (qsc *QueryStreamCreate) Start(ctx context.Context, configPk uuid.UUID) error {
	qsc.ctx = ctx
	qsc.configPK = configPk
	qsc.catch = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(qsc.errPipe)))
	if err := qsc.validateStart(); err != nil {
		return err
	}
	go qsc.listen()
	return nil
}

func (qsc *QueryStreamCreate) Send(startTS telem.TimeStamp, data *telem.ChunkData) {
	qsc.stream <- queryStreamCreateArgs{startTS: startTS, data: data}
}

func (qsc *QueryStreamCreate) Close() {
	close(qsc.stream)
	<-qsc.donePipe
	close(qsc.errPipe)
}

func (qsc *QueryStreamCreate) Errors() chan error {
	return qsc.errPipe
}

// |||| PROCESS ||||

func (qsc *QueryStreamCreate) listen() {
	qsc.updateConfigState(models.ChannelStateActive)
	defer func() {
		qsc.updateConfigState(models.ChannelStateInactive)
		qsc.donePipe <- true
	}()
	for args := range qsc.stream {
		qsc.processNextChunk(args.startTS, args.data)
	}
}

func (qsc *QueryStreamCreate) processNextChunk(startTS telem.TimeStamp, data *telem.ChunkData) {
	nextChunk := telem.NewChunk(startTS, qsc.config().DataType, qsc.config().DataRate, data)
	qsc.validateResolveNextChunk(nextChunk)

	cc := &models.ChannelChunk{
		ID:              uuid.New(),
		ChannelConfigID: qsc.config().ID,
		StartTS:         nextChunk.Start(),
		Size:            nextChunk.Size(),
	}
	ccr := &models.ChannelChunkReplica{ID: uuid.New(), ChannelChunkID: cc.ID, Telem: data}

	// CLARIFICATION: This means we tried to write a duplicate or consumed chunk.
	if cc.Size == 0 {
		return
	}

	alloc := qsc.rngSvc.NewAllocate()
	qsc.catch.Exec(alloc.Chunk(qsc.config().NodeID, cc).Exec)
	qsc.catch.Exec(alloc.ChunkReplica(ccr).Exec)

	qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(cc).Exec)
	qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(ccr).Exec)

	qsc.setPrevChunk(nextChunk)
	qsc.catch.Reset()
}

// ||| VALUE ACCESS |||

func (qsc *QueryStreamCreate) config() *models.ChannelConfig {
	if model.NewPK(qsc._config.ID).IsZero() {
		qsc.catch.Exec(query.NewRetrieve().BindExec(qsc.exec).Model(qsc._config).WherePK(qsc.configPK).Exec)
	}
	return qsc._config
}

func (qsc *QueryStreamCreate) updateConfigState(state models.ChannelState) {
	qsc.obs.Add(ObservedChannelConfig{State: state, PK: qsc.configPK})
	qsc._config.State = state
	qsc.catch.Exec(query.NewUpdate().BindExec(qsc.exec).Model(qsc._config).WherePK(qsc.configPK).Fields("State").Exec)
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

// |||| VALIDATE + RESOLVE ||||

func (qsc *QueryStreamCreate) validateStart() error {
	return validateStart().Exec(ValidateStartContext{cfg: qsc.config(), obs: qsc.obs}).Error()
}

func (qsc *QueryStreamCreate) validateResolveNextChunk(nextChunk *telem.Chunk) {
	qsc.catch.CatchSimple.Exec(func() error {
		nCtx := NextChunkContext{config: qsc.config(), prevChunk: qsc.prevChunk(), nextChunk: nextChunk}
		for _, vErr := range validateNextChunk().Exec(nCtx).Errors() {
			if rErr := qsc.resolveNextChunkError(vErr, nCtx); rErr != nil {
				return rErr
			}
		}
		return nil
	})
}

func (qsc *QueryStreamCreate) resolveNextChunkError(err error, nCtx NextChunkContext) error {
	return resolveNextChunk().Exec(err, nCtx).Error()
}

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

type NextChunkContext struct {
	ctx       context.Context
	config    *models.ChannelConfig
	prevChunk *telem.Chunk
	nextChunk *telem.Chunk
}

func validateNextChunk() *validate.Validate[NextChunkContext] {
	actions := []func(vCtx NextChunkContext) error{validateTiming}
	return validate.New[NextChunkContext](actions, errutil.WithAggregation())
}

func resolveNextChunk() *validate.Resolve[NextChunkContext] {
	actions := []func(sErr error, rCtx NextChunkContext) (bool, error){resolveTiming}
	return validate.NewResolve[NextChunkContext](actions, errutil.WithAggregation())
}
