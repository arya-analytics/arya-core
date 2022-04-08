package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	wg := &errgroup.Group{}
	fReq, rErr := sc.Receive()
	if err, done := query.StreamDone(sc.Context(), rErr); done {
		return err
	}
	sc.chunkStream = make(chan chanchunk.StreamCreateArgs)
	ctx, cancel := context.WithCancel(sc.Context())
	stream, qErr := sc.svc.NewTSCreate().
		WhereConfigPK(fReq.ConfigPK).
		Model(&sc.chunkStream).
		BindExec(sc.svc.Exec).
		Stream(ctx)
	if qErr != nil {
		cancel()
		return qErr
	}
	sc.chunkStream <- chanchunk.StreamCreateArgs{
		Start: fReq.StartTS,
		Data:  fReq.ChunkData,
	}
	sc.qStream = stream
	wg.Go(sc.relayErrors)
	wg.Go(sc.relayRequests)
	cancel()
	stream.Wait()
	log.Info("Waiting DOne")
	return wg.Wait()
}

func (sc *streamCreate) relayRequests() error {
	for {
		req, rErr := sc.Receive()
		if err, done := query.StreamDone(sc.Context(), rErr); done {
			return err
		}
		sc.chunkStream <- chanchunk.StreamCreateArgs{
			Start: req.StartTS,
			Data:  req.ChunkData,
		}
	}
}

func (sc *streamCreate) relayErrors() error {
	for err := range sc.qStream.Errors {
		if sErr, done := query.StreamDone(sc.Context(), sc.Send(StreamCreateResponse{Error: err})); done {
			return sErr
		}
	}
	return nil
}
