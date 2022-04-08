package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type StreamCreateProtocol interface {
	Context() context.Context
	Receive() (StreamCreateRequest, error)
	Send(StreamCreateResponse) error
	CloseSend()
}

type StreamCreateRequest struct {
	ConfigPK  uuid.UUID
	StartTS   telem.TimeStamp
	ChunkData *telem.ChunkData
}

type StreamCreateResponse struct {
	Error error
}

type streamCreate struct {
	StreamCreateProtocol
	svc         *chanchunk.Service
	qStream     *streamq.Stream
	chunkStream chan chanchunk.StreamCreateArgs
	cancel      context.CancelFunc
	ctx         context.Context
}

func CreateStream(svc *chanchunk.Service, cp StreamCreateProtocol) error {
	sc := &streamCreate{
		svc:                  svc,
		StreamCreateProtocol: cp,
		qStream:              &streamq.Stream{Errors: make(chan error)},
	}
	return sc.Stream()
}

func (sc *streamCreate) Stream() error {
	if err := sc.startStream(); err != nil {
		return err
	}
	wg := errgroup.Group{}
	wg.Go(sc.relayErrors)
	wg.Go(sc.relayRequests)
	return wg.Wait()
}

func (sc *streamCreate) startStream() error {
	fReq, rErr := sc.Receive()
	if err, done := query.StreamDone(rErr); done || route.CtxDone(sc.Context()) {
		return err
	}
	sc.chunkStream = make(chan chanchunk.StreamCreateArgs)
	sc.ctx, sc.cancel = context.WithCancel(sc.Context())
	stream, qErr := sc.svc.NewTSCreate().WhereConfigPK(fReq.ConfigPK).Model(&sc.chunkStream).Stream(sc.ctx)
	if qErr != nil {
		sc.cancel()
		return qErr
	}
	sc.chunkStream <- chanchunk.StreamCreateArgs{Start: fReq.StartTS, Data: fReq.ChunkData}
	sc.qStream = stream
	return nil
}

func (sc *streamCreate) relayRequests() error {
	defer sc.cancel()
	return query.StreamFor(sc.Context(), sc.Receive, func(req StreamCreateRequest) error {
		sc.chunkStream <- chanchunk.StreamCreateArgs{Start: req.StartTS, Data: req.ChunkData}
		return nil
	})
}

func (sc *streamCreate) relayErrors() error {
	return query.StreamRange[error](sc.ctx, sc.qStream.Errors, func(err error) error {
		return sc.Send(StreamCreateResponse{Error: err})
	})
}
