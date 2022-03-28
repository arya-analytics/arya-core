package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	qExec query.Execute
	delta *delta
}

func NewService(qExec query.Execute) *Service {
	d := &delta{
		inlet:     &deltaInlet{ctx: context.Background(), qExec: qExec},
		addOutlet: make(chan deltaOutlet, 1),
	}
	go d.start()
	return &Service{qExec: qExec, delta: d}
}

func (s *Service) NewStreamCreate() *StreamCreate {
	return newStreamCreate(s.qExec)
}

func (s *Service) NewStreamRetrieve() *StreamRetrieve {
	return newStreamRetrieve(s.delta)
}
