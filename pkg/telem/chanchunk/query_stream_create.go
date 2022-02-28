package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/google/uuid"
)

type QueryStreamCreateArgs struct {
	startTS telem.TimeStamp
	data    *telem.ChunkData
}

type QueryStreamCreate struct {
	config      *models.ChannelConfig
	prevChunk   *telem.Chunk
	prevChunkPK uuid.UUID
	rngSvc      *rng.Service
	errChan     chan error
	stream      chan QueryStreamCreateArgs
	ctx         context.Context
}

func (qsc *QueryStreamCreate) Start(ctx context.Context, pk uuid.UUID) {
	qsc.ctx = ctx
	qsc.listen()
}

func (qsc *QueryStreamCreate) Send(startTS telem.TimeStamp, data *telem.ChunkData) {
	qsc.stream <- QueryStreamCreateArgs{startTS: startTS, data: data}
}

func (qsc *QueryStreamCreate) Errors() chan error {
	return qsc.errChan
}

func (qsc *QueryStreamCreate) listen() {
	for args := range qsc.stream {
		alloc := qsc.rngSvc.NewAllocate()
		cc := &models.ChannelChunk{ID: uuid.New(), ChannelConfigID: qsc.config.ID, StartTS: args.startTS, Size: args.data.Size()}
		alloc.Chunk(qsc.config.NodeID, cc)
		ccr := &models.ChannelChunkReplica{ID: uuid.New(), ChannelChunkID: cc.ID, Telem: args.data}
		alloc.ChunkReplica(ccr)
		nextChunk := telem.NewChunk(args.startTS, qsc.config.DataType, qsc.config.DataRate, args.data)
		qsc.validateNextChunk(nextChunk)

	}
}

func (qsc *QueryStreamCreate) validateNextChunk(nextChunk *telem.Chunk) {
	vCtx := CreateValidateContext{prevChunk: qsc.prevChunk, nextChunk: nextChunk}
	v := createValidate().Exec(vCtx)
	if v.Error() != nil {
		v.Errors()
	}
}

type CreateValidateContext struct {
	prevChunk *telem.Chunk
	nextChunk *telem.Chunk
}

type CreateResolveContext struct {
	config *models.ChannelConfig
	*CreateValidateContext
}

func createValidate() *validate.Validate[CreateValidateContext] {
	actions := []func(vCtx CreateValidateContext) error{validateTiming}
	return validate.New[CreateValidateContext](actions)
}
