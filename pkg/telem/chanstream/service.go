package chanstream

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

type Service struct {
	qExec query.Execute
	rel   *relay
}

func NewService(qExec query.Execute) *Service {
	rel := newRelay(telem.DataRate(100), qExec)
	go rel.start()
	return &Service{qExec: qExec, rel: rel}
}

func (s *Service) NewStreamCreate() *StreamCreate {
	return newStreamCreate(s.qExec)
}

func (s *Service) NewStreamRetrieve() *StreamRetrieve {
	return newStreamRetrieve(s.rel)
}
