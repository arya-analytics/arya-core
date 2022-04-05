// Package chanchunk provides an interface for reading and modifying bulk telemetry on the cluster.
//
// chanchunk.Service is the key interface for this package.
//
package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// Service provides an interface for reading and modifying bulk telemetry on the cluster.
// Avoid constructing directly, and instead call NewService.
//
// Operations:
//		NewStreamCreate -> Save contiguous chunks of telemetry.
//
type Service struct {
	obs    observe
	qa     query.Assemble
	rngSvc *rng.Service
}

// NewService creates a new service with the provided parameters.
func NewService(qa query.Assemble, rngSVC *rng.Service) *Service {
	return &Service{qa: qa, obs: newObserveMem(), rngSvc: rngSVC}
}

// NewStreamCreate opens a StreamCreate.
func (s *Service) NewStreamCreate() *StreamCreate {
	return newStreamCreate(s.qa.Exec, s.obs, s.rngSvc)
}

// NewStreamRetrieve opens a StreamCreate.
func (s *Service) NewStreamRetrieve() *StreamRetrieve {
	return newStreamRetrieve(s.qa)
}
