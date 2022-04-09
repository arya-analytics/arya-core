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
	// Receive the first request as it has the config PK.
	fReq, err := sc.Receive()
	// If our stream already errored out or was cancelled, return the error.
	if dErr, done := query.StreamDone(sc.Context(), err); done {
		return dErr
	}
	sc.chunkStream = make(chan chanchunk.StreamCreateArgs)
	// Creates a context that we can use to cancel the query at the end of the stream.
	sc.ctx, sc.cancel = context.WithCancel(sc.Context())
	// If our query had bad parameters or encountered issues during assembly, return the error.
	if sc.qStream, err = sc.svc.NewTSCreate().WhereConfigPK(fReq.ConfigPK).Model(&sc.chunkStream).Stream(sc.ctx); err != nil {
		sc.cancel()
		return err
	}
	// We still need to pipe the first chunk to be created.
	sc.chunkStream <- chanchunk.StreamCreateArgs{Start: fReq.StartTS, Data: fReq.ChunkData}
	return nil
}

func (sc *streamCreate) relayRequests() error {
	// Cancel the context to break out of the error relay and tell the query to stop.
	defer sc.cancel()
	// Relay all requests to the query.
	return query.StreamFor(sc.ctx, sc.Receive, func(req StreamCreateRequest) error {
		sc.chunkStream <- chanchunk.StreamCreateArgs{Start: req.StartTS, Data: req.ChunkData}
		return nil
	})
}

func (sc *streamCreate) relayErrors() error {
	return query.StreamRange(sc.ctx, sc.qStream.Errors, func(err error) error {
		log.Fatal(err)
		return sc.Send(StreamCreateResponse{Error: err})
	})
}
