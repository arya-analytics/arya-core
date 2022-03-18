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
	exec   query.Execute
	rngSVC *rng.Service
}

// NewService creates a new service with the provided parameters.
func NewService(exec query.Execute, rngSVC *rng.Service) *Service {
	return &Service{exec: exec, obs: newObserveMem(), rngSVC: rngSVC}
}

// NewStreamCreate opens a QueryStreamCreate.
func (s *Service) NewStreamCreate() *QueryStreamCreate {
	return newStreamCreate(s.exec, s.obs, s.rngSVC)
}
