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

const errorPipeCapacity = 10

// StreamCreate creates a set of contiguous data chunks.
// Avoid instantiating directly, and instead instantiate by calling Service.NewStreamCreate.
type StreamCreate struct {
	obs        observe
	exec       query.Execute
	rngSvc     *rng.Service
	configPK   uuid.UUID
	_config    *models.ChannelConfig
	_prevChunk *telem.Chunk
	prevCCPK   uuid.UUID
	errors     chan error
	catch      *errutil.CatchContext
	done       chan bool
	stream     chan streamCreateArgs
	ctx        context.Context
}

type streamCreateArgs struct {
	start telem.TimeStamp
	data  *telem.ChunkData
}

func newStreamCreate(qExec query.Execute, obs observe, rngSvc *rng.Service) *StreamCreate {
	return &StreamCreate{
		obs:     obs,
		exec:    qExec,
		rngSvc:  rngSvc,
		_config: &models.ChannelConfig{},
		errors:  make(chan error, errorPipeCapacity),
		stream:  make(chan streamCreateArgs),
		done:    make(chan bool),
	}
}

// Start starts stream. Start must be called before Send. Returns any errors encountered during stream start.
func (qsc *StreamCreate) Start(ctx context.Context, configPk uuid.UUID) error {
	qsc.ctx = ctx
	qsc.configPK = configPk
	qsc.catch = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(qsc.errors)))
	if err := qsc.validateStart(); err != nil {
		return err
	}
	go qsc.listen()
	return nil
}

// Send creates a new chunk of data starting at the specified timestamp.
func (qsc *StreamCreate) Send(start telem.TimeStamp, data *telem.ChunkData) {
	qsc.stream <- streamCreateArgs{start: start, data: data}
}

// Close safely closes the stream.
func (qsc *StreamCreate) Close() {
	close(qsc.stream)
	<-qsc.done
	close(qsc.errors)
}

// Errors returns errors encountered during stream operation.
func (qsc *StreamCreate) Errors() chan error {
	return qsc.errors
}

// |||| PROCESS ||||

func (qsc *StreamCreate) listen() {
	qsc.updateConfigStatus(models.ChannelStatusActive)
	defer func() {
		qsc.updateConfigStatus(models.ChannelStatusInactive)
		qsc.done <- true
	}()
	for args := range qsc.stream {
		qsc.processNextChunk(args.start, args.data)
	}
}

func (qsc *StreamCreate) processNextChunk(startTS telem.TimeStamp, data *telem.ChunkData) {
	nc := telem.NewChunk(startTS, qsc.config().DataType, qsc.config().DataRate, data)
	qsc.validateResolveNextChunk(nc)

	cc := &models.ChannelChunk{
		ID:              uuid.New(),
		ChannelConfigID: qsc.config().ID,
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

	a := qsc.rngSvc.NewAllocate()
	qsc.catch.Exec(a.Chunk(qsc.config().NodeID, cc).Exec)
	qsc.catch.Exec(a.ChunkReplica(ccr).Exec)

	qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(cc).Exec)
	qsc.catch.Exec(query.NewCreate().BindExec(qsc.exec).Model(ccr).Exec)

	qsc.setPrevChunk(nc)
	qsc.catch.Reset()
}

// ||| VALUE ACCESS |||

func (qsc *StreamCreate) config() *models.ChannelConfig {
	if model.NewPK(qsc._config.ID).IsZero() {
		qsc.catch.Exec(query.NewRetrieve().BindExec(qsc.exec).Model(qsc._config).WherePK(qsc.configPK).Exec)
	}
	return qsc._config
}

func (qsc *StreamCreate) updateConfigStatus(status models.ChannelStatus) {
	qsc.obs.Add(observedChannelConfig{Status: status, PK: qsc.configPK})
	qsc._config.Status = status
	qsc.catch.CatchSimple.Exec(func() error {
		return query.NewUpdate().
			BindExec(qsc.exec).
			Model(qsc._config).
			WherePK(qsc.configPK).
			Fields("Status").Exec(context.Background())
	})
}

func (qsc *StreamCreate) prevChunk() *telem.Chunk {
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

func (qsc *StreamCreate) setPrevChunk(chunk *telem.Chunk) {
	qsc._prevChunk = chunk
}

// |||| VALIDATE + RESOLVE ||||

func (qsc *StreamCreate) validateStart() error {
	return validateStart().Exec(validateStartContext{cfg: qsc.config(), obs: qsc.obs}).Error()
}

func (qsc *StreamCreate) validateResolveNextChunk(nextChunk *telem.Chunk) {
	nc := nextChunkContext{cfg: qsc.config(), prev: qsc.prevChunk(), next: nextChunk}
	qsc.catch.CatchSimple.Exec(func() error {
		for _, vErr := range validateNextChunk().Exec(nc).Errors() {
			if rErr := qsc.resolveNextChunkError(vErr, nc); rErr != nil {
				return rErr
			}
		}
		return nil
	})
}

func (qsc *StreamCreate) resolveNextChunkError(err error, nCtx nextChunkContext) error {
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
		return errors.New("open a second stream to an active channel")
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
