// Package rng holds utilities for managing and working with models.Range, models.RangeReplica, and models.RangeLease
// objects. This includes range partitioning algorithms (PartitionExecute), allocation utilities (Allocate), and, in
// the future, range replication services.
package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

// Service is the central access point to the rng package. Provides utilities for allocating channel chunks to range
// as well as starting and stopping rng specific tasks, such as partitioning rngMap.
// ONLY one Service should exist per core instance.
type Service struct {
	ps   tasks.Schedule
	obs  Observe
	exec query.Execute
}

// NewService creates a new rng.Service. Requires a val.
func NewService(obs Observe, exec query.Execute) *Service {
	return &Service{obs: obs, exec: exec}
}

// NewAllocate creates a new Allocate and returns it.
func (s *Service) NewAllocate() *Allocate {
	return &Allocate{obs: s.obs, pst: s.pst}
}

// Start starts Service internal tasks.
// NOTE: If restarting the Service, call Stop before calling Start again.
func (s *Service) Start(ctx context.Context, opts ...tasks.ScheduleOpt) {
	s.ps = newSchedulePartition(&partitionDetect{Persist: s.pst, Observe: s.obs}, opts...)
	go s.ps.Start(ctx)
}

// Stop stops Service internal tasks.
func (s *Service) Stop() {
	s.ps.Stop()
	s.ps = nil
}

// Errors returns a channel that sends any errors encountered during Service operation.
func (s *Service) Errors() chan error {
	return s.ps.Errors()
}
