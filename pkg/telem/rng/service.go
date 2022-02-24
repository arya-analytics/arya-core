package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

type Service struct {
	ps  tasks.Schedule
	obs Observe
	pst Persist
}

func NewService(obs Observe, p Persist) *Service {
	return &Service{obs: obs, pst: p}
}

func (s *Service) NewAllocate() *Allocate {
	return &Allocate{obs: s.obs, pst: s.pst}
}

func (s *Service) Start(ctx context.Context, opts ...tasks.ScheduleOpt) {
	s.ps = newSchedulerPartition(&PartitionDetect{Persist: s.pst, Observe: s.obs}, opts...)
	s.ps.Start(ctx)
}

func (s *Service) Stop() {
	s.ps.Stop()
	s.ps = nil
}

func (s *Service) Errors() chan error {
	return s.ps.Errors()
}
