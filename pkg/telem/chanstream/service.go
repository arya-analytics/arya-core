package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type Service struct {
	qExec query.Execute
	delta *route.Delta[*models.ChannelSample, outletContext]
}

func NewService(qExec query.Execute) *Service {
	d := route.NewDelta[*models.ChannelSample, outletContext](&deltaInlet{ctx: context.Background(), qExec: qExec})
	go d.Start()
	return &Service{qExec: qExec, delta: d}
}

func (s *Service) NewStreamCreate() *StreamCreate {
	return newStreamCreate(s.qExec)
}

func (s *Service) NewStreamRetrieve() *StreamRetrieve {
	return newStreamRetrieve(s.delta)
}
