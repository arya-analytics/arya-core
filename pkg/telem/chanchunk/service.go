package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	obs    Observe
	exec   query.Execute
	rngSVC *rng.Service
}

func NewService(exec query.Execute, obs Observe, rngSVC *rng.Service) *Service {
	return &Service{exec: exec, obs: obs, rngSVC: rngSVC}
}

func (s *Service) NewStreamCreate() *QueryStreamCreate {
	return newStreamCreate(s.exec, s.obs, s.rngSVC)
}
