package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"golang.org/x/sync/errgroup"
)

type CreateProtocol interface {
	Context() context.Context
	Receive() (CreateRequest, error)
	Send(CreateResponse) error
}

type CreateRequest struct {
	Sample *models.ChannelSample
}

type CreateResponse struct {
	Error error
}

type Create struct {
	CreateProtocol
	svc          *chanstream.Service
	qStream      *streamq.Stream
	sampleStream chan *models.ChannelSample
}

func CreateStream(svc *chanstream.Service, rp CreateProtocol) error {
	c := &Create{
		CreateProtocol: rp,
		svc:            svc,
	}
	return c.Stream()
}

func (c *Create) Stream() error {
	wg := errgroup.Group{}
	c.sampleStream = make(chan *models.ChannelSample)
	stream, qErr := streamq.NewTSCreate().Model(&c.sampleStream).BindExec(c.svc.Exec).Stream(c.Context())
	if qErr != nil {
		return qErr
	}
	c.qStream = stream
	wg.Go(c.relayErrors)
	wg.Go(c.relayRequests)
	return wg.Wait()
}

func (c *Create) relayErrors() error {
	for err := range c.qStream.Errors {
		if sErr, done := query.StreamDone(c.Context(), c.Send(CreateResponse{Error: err})); done {
			return sErr
		}
	}
	return nil
}

func (c *Create) relayRequests() error {
	for {
		req, rErr := c.Receive()
		if err, done := query.StreamDone(c.Context(), rErr); done {
			return err
		}
		c.sampleStream <- req.Sample
	}
}
