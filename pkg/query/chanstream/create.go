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

type create struct {
	CreateProtocol
	svc          *chanstream.Service
	qStream      *streamq.Stream
	sampleStream chan *models.ChannelSample
	ctx          context.Context
	cancelQ      context.CancelFunc
}

func CreateStream(svc *chanstream.Service, rp CreateProtocol) error {
	c := &create{CreateProtocol: rp, svc: svc}
	return c.Stream()
}

func (c *create) Stream() error {
	if err := c.startStream(); err != nil {
		return err
	}
	wg := errgroup.Group{}
	wg.Go(c.relayErrors)
	wg.Go(c.relayRequests)
	return wg.Wait()
}

func (c *create) startStream() (err error) {
	c.sampleStream = make(chan *models.ChannelSample)
	c.ctx, c.cancelQ = context.WithCancel(c.Context())
	c.qStream, err = streamq.NewTSCreate().Model(&c.sampleStream).BindExec(c.svc.Exec).Stream(c.Context())
	return err
}

func (c *create) relayErrors() error {
	return query.StreamRange(c.ctx, c.qStream.Errors, func(err error) error {
		return c.Send(CreateResponse{Error: err})
	})
}

func (c *create) relayRequests() error {
	defer c.cancelQ()
	return query.StreamFor(c.ctx, c.Receive, func(req CreateRequest) error {
		c.sampleStream <- req.Sample
		return nil
	})
}
