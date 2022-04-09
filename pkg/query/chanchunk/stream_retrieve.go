package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type StreamRetrieveProtocol interface {
	Context() context.Context
	Send(StreamRetrieveResponse) error
}

type StreamRetrieveRequest struct {
	ChannelConfigID uuid.UUID
	TimeRange       telem.TimeRange
}

type StreamRetrieveResponse struct {
	StartTS  telem.TimeStamp
	DataType telem.DataType
	DataRate telem.DataRate
	Data     *telem.ChunkData
	Error    error
}

type streamRetrieve struct {
	StreamRetrieveProtocol
	svc         *chanchunk.Service
	qStream     *streamq.Stream
	chunkStream chan *telem.Chunk
	ctx         context.Context
	cancel      context.CancelFunc
}

func RetrieveStream(svc *chanchunk.Service, sp StreamRetrieveProtocol, req StreamRetrieveRequest) error {
	sr := &streamRetrieve{
		svc:                    svc,
		StreamRetrieveProtocol: sp,
		chunkStream:            make(chan *telem.Chunk),
	}
	return sr.Stream(req)
}

func (sr *streamRetrieve) Stream(req StreamRetrieveRequest) error {
	if err := sr.startStream(req); err != nil {
		return err
	}
	wg := &errgroup.Group{}
	wg.Go(sr.relayErrors)
	wg.Go(sr.relayChunks)
	return wg.Wait()
}

func (sr *streamRetrieve) startStream(req StreamRetrieveRequest) (err error) {
	sr.chunkStream = make(chan *telem.Chunk)
	sr.ctx, sr.cancel = context.WithCancel(sr.Context())
	sr.qStream, err = sr.svc.NewTSRetrieve().
		WhereConfigPK(req.ChannelConfigID).
		Model(&sr.chunkStream).
		WhereTimeRange(req.TimeRange).
		Stream(sr.ctx)
	if err != nil {
		sr.cancel()
	}
	return err
}

func (sr *streamRetrieve) relayErrors() error {
	return query.StreamRange[error](sr.ctx, sr.qStream.Errors, func(err error) error {
		return sr.Send(StreamRetrieveResponse{Error: err})
	})
}

func (sr *streamRetrieve) relayChunks() error {
	ctx, cancel := context.WithCancel(sr.ctx)
	defer func() {
		// Make sure the query is cancelled.
		cancel()
		// Wait for the query to complete.
		sr.qStream.Wait()
		// Cancel error streaming.
		sr.cancel()
	}()
	return query.StreamRange(ctx, sr.chunkStream, func(chunk *telem.Chunk) error {
		return sr.Send(StreamRetrieveResponse{
			StartTS:  chunk.Start(),
			DataType: chunk.DataType,
			DataRate: chunk.DataRate,
			Data:     chunk.ChunkData,
		})
	})
}
