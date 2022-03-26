package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

type StreamCreate struct {
	ctx    context.Context
	qExec  query.Execute
	stream chan *models.ChannelSample
	ge     tsquery.GoExecOpt
	errors chan error
}

const errorPipeCapacity = 10

func newStreamCreate(qExec query.Execute) *StreamCreate {
	return &StreamCreate{
		errors: make(chan error, errorPipeCapacity),
		stream: make(chan *models.ChannelSample, 1),
		qExec:  qExec,
		ge:     tsquery.GoExecOpt{Errors: make(chan error, errorPipeCapacity)},
	}
}

func (sc *StreamCreate) Start(ctx context.Context) {
	sc.ctx = ctx
	sc.listen()
}

func (sc *StreamCreate) listen() {
	sc.ge = tsquery.NewCreate().Model(&sc.stream).BindExec(sc.qExec).GoExec(sc.ctx)
}

func (sc *StreamCreate) Send(sample *models.ChannelSample) {
	sc.stream <- sample
}

func (sc *StreamCreate) Stop() {
	sc.ge.Release()
}

func (sc *StreamCreate) Errors() chan error {
	return sc.ge.Errors
}
