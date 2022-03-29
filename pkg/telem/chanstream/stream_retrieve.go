package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type StreamRetrieve struct {
	delta  *route.Delta[*models.ChannelSample, outletContext]
	stream chan *models.ChannelSample
	errors chan error
	pkc    model.PKChain
}

func newStreamRetrieve(delta *route.Delta[*models.ChannelSample, outletContext]) *StreamRetrieve {
	return &StreamRetrieve{
		delta:  delta,
		stream: make(chan *models.ChannelSample),
		errors: make(chan error, errorPipeCapacity),
	}
}

// || QUERY ||

func (s *StreamRetrieve) Start(ctx context.Context) chan *models.ChannelSample {
	s.delta.AddOutlet(s)
	return s.stream
}

func (s *StreamRetrieve) WherePKC(pkc model.PKChain) *StreamRetrieve {
	s.pkc = pkc
	s.stream = make(chan *models.ChannelSample, len(s.pkc))
	return s
}

func (s *StreamRetrieve) Close() {
	s.delta.RemoveOutlet(s)
	close(s.stream)
	close(s.errors)
}

// || DELTA OUTLET IMPL ||

func (s *StreamRetrieve) Errors() chan error {
	return s.errors
}

func (s *StreamRetrieve) SendError(e error) {
	s.errors <- e
}

func (s *StreamRetrieve) Send(v *models.ChannelSample) {
	s.stream <- v
}

func (s *StreamRetrieve) Context() outletContext {
	return outletContext{pkc: s.pkc}
}

type outletContext struct {
	pkc model.PKChain
}

// |||| INLET ||||

type deltaInlet struct {
	ctx   context.Context
	qExec query.Execute
	s     chan *models.ChannelSample
	goe   tsquery.GoExecOpt
}

func (i *deltaInlet) Stream() chan *models.ChannelSample {
	return i.s
}

func (i *deltaInlet) Errors() chan error {
	return i.goe.Errors
}

func (i *deltaInlet) Update(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) {
	var allPKC model.PKChain
	for o := range dCtx.Outlets {
		allPKC = append(allPKC, o.Context().pkc...)
	}
	allPKC = allPKC.Unique()
	i.s = make(chan *models.ChannelSample, len(allPKC))
	i.goe = tsquery.NewRetrieve().Model(&i.s).BindExec(i.qExec).WherePKs(allPKC).GoExec(i.ctx)
}
