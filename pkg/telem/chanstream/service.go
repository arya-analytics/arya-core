package chanstream

import "github.com/arya-analytics/aryacore/pkg/util/query"

type Service struct {
	qExec query.Execute
}

func NewService(qExec query.Execute) *Service {
	return &Service{qExec: qExec}
}

func (s *Service) NewStreamCreate() *StreamCreate {
	return newStreamCreate(s.qExec)
}
