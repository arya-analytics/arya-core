package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"golang.org/x/sync/errgroup"
)

type RetrieveProtocol interface {
	Context() context.Context
	Receive() (RetrieveRequest, error)
	Send(RetrieveResponse) error
}

type RetrieveRequest struct {
	PKC model.PKChain
}

type RetrieveResponse struct {
	Sample *models.ChannelSample
	Error  error
}

func RetrieveStream(svc *chanstream.Service, rp RetrieveProtocol) error {
	r := &retrieve{
		RetrieveProtocol: rp,
		svc:              svc,
		qStream:          &streamq.Stream{Errors: make(chan error, 10)},
		sampleStream:     make(chan *models.ChannelSample),
		updateSig:        make(chan struct{}),
	}
	return r.stream()
}

type retrieve struct {
	RetrieveProtocol
	_cancelQuery context.CancelFunc
	cancelRelay  context.CancelFunc
	relayCtx     context.Context
	svc          *chanstream.Service
	qStream      *streamq.Stream
	sampleStream chan *models.ChannelSample
	updateSig    chan struct{}
}

func (r *retrieve) stream() error {
	r.relayCtx, r.cancelRelay = context.WithCancel(r.Context())
	wg := errgroup.Group{}
	wg.Go(r.relayErrors)
	wg.Go(r.relaySamples)
	wg.Go(r.listenForUpdates)
	return wg.Wait()
}

func (r *retrieve) relayErrors() error {
	return query.StreamRange(r.relayCtx, r.qStream.Errors, func(err error) error {
		return r.Send(RetrieveResponse{Error: err})
	})
}

func (r *retrieve) relaySamples() error {
	for {
		select {
		case s := <-r.sampleStream:
			if err, done := query.StreamDone(r.relayCtx, r.Send(RetrieveResponse{Sample: s})); done {
				return err
			}
		// When we receive an update signal, it means we've changed the value of r.sampleStream and need to
		// restart the loop.
		case <-r.updateSig:
			continue
		}
	}
}

func (r *retrieve) listenForUpdates() error {
	defer r.cancelQuery()
	defer r.cancelRelay()
	return query.StreamFor(r.Context(), r.Receive, func(req RetrieveRequest) error {
		r.updateQuery(req.PKC)
		return nil
	})
}

func (r *retrieve) updateQuery(pkc model.PKChain) {
	pSampleStream := make(chan *models.ChannelSample, len(pkc))
	ctx, cancel := context.WithCancel(context.Background())
	pqStream, err := r.svc.NewTSRetrieve().Model(&pSampleStream).WherePKs(pkc).Stream(ctx)
	if err != nil {
		cancel()
		r.qStream.Errors <- err
	}
	r.cancelQuery()
	r.sampleStream = pSampleStream
	r._cancelQuery = cancel
	r.qStream = pqStream
	r.updateSig <- struct{}{}
}

func (r *retrieve) cancelQuery() {
	if r._cancelQuery != nil {
		r._cancelQuery()
	}
}
