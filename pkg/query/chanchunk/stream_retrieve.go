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
	wg := &errgroup.Group{}
	sr.chunkStream = make(chan *telem.Chunk)
	stream, qErr := sr.svc.NewTSRetrieve().
		WhereConfigPK(req.ChannelConfigID).
		Model(&sr.chunkStream).
		WhereTimeRange(req.TimeRange).
		Stream(sr.Context())
	if qErr != nil {
		return qErr
	}
	sr.qStream = stream
	wg.Go(sr.relayErrors)
	wg.Go(sr.relayChunks)
	return wg.Wait()
}

func (sr *streamRetrieve) relayErrors() error {
	for err := range sr.qStream.Errors {
		if sErr, done := query.StreamDone(sr.Context(), sr.Send(StreamRetrieveResponse{Error: err})); done {
			return sErr
		}
	}
	return nil
}

func (sr *streamRetrieve) relayChunks() error {
	for chunk := range sr.chunkStream {
		if sErr, done := query.StreamDone(sr.Context(), sr.Send(StreamRetrieveResponse{
			StartTS:  chunk.Start(),
			DataType: chunk.DataType,
			DataRate: chunk.DataRate,
			Data:     chunk.ChunkData,
		})); done {
			return sErr
		}
	}
	return nil
}
