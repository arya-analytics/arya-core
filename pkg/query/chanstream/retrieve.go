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

type Retrieve struct {
	RetrieveProtocol
	cancelQ      context.CancelFunc
	svc          *chanstream.Service
	qStream      *streamq.Stream
	sampleStream chan *models.ChannelSample
	updateSig    chan struct{}
}

func RetrieveStream(svc *chanstream.Service, rp RetrieveProtocol) error {
	r := &Retrieve{
		RetrieveProtocol: rp,
		svc:              svc,
		qStream:          &streamq.Stream{Errors: make(chan error, 10)},
		sampleStream:     make(chan *models.ChannelSample),
		updateSig:        make(chan struct{}),
	}
	return r.Stream()
}

func (r *Retrieve) Stream() error {
	defer r.closeQuery()
	wg := errgroup.Group{}
	wg.Go(r.relayErrors)
	wg.Go(r.relaySamples)
	wg.Go(r.listenForUpdates)
	return wg.Wait()
}

func (r *Retrieve) relayErrors() error {
	for err := range r.qStream.Errors {
		if sErr, done := query.StreamDone(r.Context(), r.Send(RetrieveResponse{Error: err})); done {
			return sErr
		}
	}
	return nil
}

func (r *Retrieve) relaySamples() error {
	for {
		select {
		case s := <-r.sampleStream:
			if err, done := query.StreamDone(r.Context(), r.Send(RetrieveResponse{Sample: s})); done {
				return err
			}
		case <-r.updateSig:
			continue
		}
	}
}

func (r *Retrieve) listenForUpdates() error {
	for {
		req, rErr := r.Receive()
		if err, done := query.StreamDone(r.Context(), rErr); done {
			return err
		}
		if uErr := r.updateQuery(req.PKC); uErr != nil {
			r.qStream.Errors <- uErr
		}
	}
}

func (r *Retrieve) updateQuery(pkc model.PKChain) error {
	pSampleStream := make(chan *models.ChannelSample, len(pkc))
	ctx, cancel := context.WithCancel(context.Background())
	pqStream, err := streamq.
		NewTSRetrieve().
		Model(&pSampleStream).
		WherePKs(pkc).
		BindExec(r.svc.Exec).
		Stream(ctx)
	if err != nil {
		cancel()
		return err
	}
	r.closeQuery()
	r.sampleStream = pSampleStream
	r.cancelQ = cancel
	r.qStream = pqStream
	r.updateSig <- struct{}{}
	return nil
}

func (r *Retrieve) closeQuery() {

}
