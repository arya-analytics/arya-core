package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
)

type QueryStreamCreateArgs struct {
	startTS telem.TimeStamp
	data    *telem.ChunkData
}

type QueryStreamCreate struct {
	cluster    cluster.Cluster
	rngSvc     *rng.Service
	configPK   uuid.UUID
	_config    *models.ChannelConfig
	_prevChunk *telem.Chunk
	prevCCPK   uuid.UUID
	errChan    chan error
	stream     chan QueryStreamCreateArgs
	ctx        context.Context
}

func newStreamCreate(cluster cluster.Cluster, rngSvc *rng.Service) *QueryStreamCreate {
	return &QueryStreamCreate{
		cluster: cluster,
		rngSvc:  rngSvc,
		_config: &models.ChannelConfig{},
		errChan: make(chan error, 10),
		stream:  make(chan QueryStreamCreateArgs),
	}
}

func (qsc *QueryStreamCreate) Start(ctx context.Context, pk uuid.UUID) *QueryStreamCreate {
	qsc.ctx, qsc.configPK = ctx, pk
	qsc.listen()
	return qsc
}

func (qsc *QueryStreamCreate) config() *models.ChannelConfig {
	if model.NewPK(qsc._config.ID).IsZero() {
		if err := qsc.cluster.NewRetrieve().Model(qsc._config).WherePK(qsc.configPK).Exec(qsc.ctx); err != nil {
			qsc.Errors() <- err
		}
	}
	return qsc._config
}

func (qsc *QueryStreamCreate) prevChunk() *telem.Chunk {
	if qsc._prevChunk == nil {
		ccr := &models.ChannelChunkReplica{}
		if err := qsc.cluster.NewRetrieve().
			Model(ccr).
			Relation("ChannelChunk", "ID", "StartTS", "Size").
			WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": qsc.config().ID}).Exec(qsc.ctx); err != nil {
			qsc.Errors() <- err
		}
		qsc._prevChunk = telem.NewChunk(ccr.ChannelChunk.StartTS, qsc.config().DataType, qsc.config().DataRate, ccr.Telem)
	}
	return qsc._prevChunk
}

func (qsc *QueryStreamCreate) Send(startTS telem.TimeStamp, data *telem.ChunkData) {
	qsc.stream <- QueryStreamCreateArgs{startTS: startTS, data: data}
}

func (qsc *QueryStreamCreate) Close() {
	close(qsc.stream)
	<-qsc.Errors()
}

func (qsc *QueryStreamCreate) Errors() chan error {
	return qsc.errChan
}

func (qsc *QueryStreamCreate) listen() {
	qsc.prevChunk()
	for args := range qsc.stream {
		c := errutil.NewCatchWCtx(qsc.ctx)
		alloc := qsc.rngSvc.NewAllocate()
		cc := &models.ChannelChunk{ID: uuid.New(), ChannelConfigID: qsc.config().ID}
		ccr := &models.ChannelChunkReplica{ID: uuid.New(), ChannelChunkID: cc.ID, Telem: args.data}

		c.Exec(alloc.Chunk(qsc.config().NodeID, cc).Exec)
		c.Exec(alloc.ChunkReplica(ccr).Exec)

		nextChunk := telem.NewChunk(args.startTS, qsc.config().DataType, qsc.config().DataRate, args.data)

		c.CatchSimple.Exec(func() error { return qsc.validateNextChunk(nextChunk) })

		cc.Size, cc.StartTS = nextChunk.Size(), nextChunk.Start()

		log.Info("made it here")
		c.Exec(qsc.cluster.NewCreate().Model(cc).Exec)
		c.Exec(qsc.cluster.NewCreate().Model(ccr).Exec)

		if c.Error() != nil {
			qsc.Errors() <- c.Error()
		}
	}
	qsc.Errors() <- io.EOF
}

func (qsc *QueryStreamCreate) validateNextChunk(nextChunk *telem.Chunk) error {
	vCtx := CreateValidateContext{prevChunk: qsc.prevChunk(), nextChunk: nextChunk}
	for _, err := range createValidate().Exec(vCtx).Errors() {
		if rErr := qsc.resolveNextChunkError(err, vCtx); rErr != nil {
			return rErr
		}
	}
	return nil
}

func (qsc *QueryStreamCreate) resolveNextChunkError(err error, vCtx CreateValidateContext) error {
	return createResolve().Exec(err, CreateResolveContext{config: qsc.config(), CreateValidateContext: vCtx}).Error()
}

// |||| VALIDATE ||||

type CreateValidateContext struct {
	ctx       context.Context
	prevChunk *telem.Chunk
	nextChunk *telem.Chunk
}

func createValidate() *validate.Validate[CreateValidateContext] {
	actions := []func(vCtx CreateValidateContext) error{validateTiming}
	return validate.New[CreateValidateContext](actions, validate.WithAggregation())
}

// |||| RESOLVE ||||

type CreateResolveContext struct {
	config *models.ChannelConfig
	CreateValidateContext
}

func createResolve() *validate.Resolve[CreateResolveContext] {
	actions := []func(sErr error, rCtx CreateResolveContext) (bool, error){resolveTiming}
	return validate.NewResolve[CreateResolveContext](actions, validate.WithAggregation())
}
