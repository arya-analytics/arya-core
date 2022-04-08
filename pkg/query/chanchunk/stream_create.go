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

type StreamCreateProtocol interface {
	Context() context.Context
	Receive() (StreamCreateRequest, error)
	Send(StreamCreateResponse) error
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
	stream, qErr := sc.svc.NewTSCreate().
		Model(&sc.chunkStream).
		BindExec(sc.svc.Exec).
		Stream(sc.Context(), chanchunk.ContextArg{ConfigPK: fReq.ConfigPK})
	if qErr != nil {
		return qErr
	}
	sc.qStream = stream
	wg.Go(sc.relayErrors)
	wg.Go(sc.relayRequests)
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
