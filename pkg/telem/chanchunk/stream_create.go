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

// Start starts streamq. Start must be called before Send. Returns any errors encountered during streamq start.
func (sc *StreamCreate) Start(ctx context.Context, configPk uuid.UUID) error {
	sc.ctx = ctx
	sc.configPK = configPk
	sc.catch = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(sc.errors)))
	if err := sc.validateStart(); err != nil {
		return err
	}
	go sc.listen()
	return nil
}

// Send creates a new chunk of data starting at the specified timestamp.
func (sc *StreamCreate) Send(start telem.TimeStamp, data *telem.ChunkData) {
	sc.stream <- streamCreateArgs{start: start, data: data}
}

// Close safely closes the streamq.
func (sc *StreamCreate) Close() {
	close(sc.stream)
	<-sc.done
	close(sc.errors)
}

// Errors returns errors encountered during streamq operation.
func (sc *StreamCreate) Errors() chan error {
	return sc.errors
}

// |||| PROCESS ||||

func (sc *StreamCreate) listen() {
	sc.updateConfigStatus(models.ChannelStatusActive)
	defer func() {
		sc.updateConfigStatus(models.ChannelStatusInactive)
		sc.done <- true
	}()
	for args := range sc.stream {
		sc.processNextChunk(args.start, args.data)
	}
}

func (sc *StreamCreate) processNextChunk(startTS telem.TimeStamp, data *telem.ChunkData) {
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

	sc.catch.Exec(query.NewCreate().BindExec(sc.exec).Model(cc).Exec)
	sc.catch.Exec(query.NewCreate().BindExec(sc.exec).Model(ccr).Exec)

	sc.setPrevChunk(nc)
	sc.catch.Reset()
}

// ||| VALUE ACCESS |||

func (sc *StreamCreate) config() *models.ChannelConfig {
	if model.NewPK(sc._config.ID).IsZero() {
		sc.catch.Exec(query.NewRetrieve().BindExec(sc.exec).Model(sc._config).WherePK(sc.configPK).Exec)
	}
	return sc._config
}

func (sc *StreamCreate) updateConfigStatus(status models.ChannelStatus) {
	sc.obs.Add(observedChannelConfig{Status: status, PK: sc.configPK})
	sc._config.Status = status
	sc.catch.CatchSimple.Exec(func() error {
		return query.NewUpdate().
			BindExec(sc.exec).
			Model(sc._config).
			WherePK(sc.configPK).
			Fields("Status").Exec(context.Background())
	})
}

func (sc *StreamCreate) prevChunk() *telem.Chunk {
	if sc._prevChunk == nil {
		sc.catch.Exec(func(ctx context.Context) error {
			ccr := &models.ChannelChunkReplica{}
			err := query.NewRetrieve().
				BindExec(sc.exec).
				Model(ccr).
				Relation("ChannelChunk", "ID", "StartTS", "Size").
				WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": sc.config().ID}).Exec(sc.ctx)
			sErr, ok := err.(query.Error)
			if !ok || sErr.Type != query.ErrorTypeItemNotFound {
				return err
			}
			// If we don't find the item, this isn't an exceptional case, it just means the channel doesn't have any
			// data, so we can just return nil early.
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
	if sCtx.cfg.Status == models.ChannelStatusActive || oc.Status == models.ChannelStatusActive {
		return errors.New("open a second streamq to an active channel")
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
