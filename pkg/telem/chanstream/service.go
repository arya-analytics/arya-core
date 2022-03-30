package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type Service struct {
	qExec query.Execute
	delta *route.Delta[*models.ChannelSample, outletContext]
}

func NewService(qExec query.Execute) *Service {
	d := route.NewDelta[*models.ChannelSample, outletContext](&deltaInlet{qExec: qExec, qStream: &streamq.Stream{Errors: make(chan error, errorPipeCapacity)}})
	go d.Start()
	return &Service{qExec: qExec, delta: d}
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSRetrieve{}: newStreamRetrieve(s.delta).exec,
		&streamq.TSCreate{}:   newStreamCreate(s.qExec).exec,
	})
}

func (s *Service) NewTSCreate() *streamq.TSCreate {
	return streamq.NewTSCreate().BindExec(s.Exec)
}

func (s *Service) NewTSRetrieve() *streamq.TSRetrieve {
	return streamq.NewTSRetrieve().BindExec(s.Exec)
}
