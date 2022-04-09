// Package chanchunk provides an interface for reading and modifying bulk telemetry on the cluster.
//
// chanchunk.Service is the key interface for this package.
//
package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
)

// Service provides an interface for reading and modifying bulk telemetry on the cluster.
// Avoid constructing directly, and instead call NewService.
type Service struct {
	obs    observe
	rngSvc *rng.Service
	qExec  query.Execute
}

// NewService creates a new service with the provided parameters.
func NewService(qExec query.Execute, rngSvc *rng.Service) *Service {
	svc := &Service{qExec: qExec, obs: newObserveMem(), rngSvc: rngSvc}
	return svc
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSRetrieve{}: newStreamRetrieve(s.qExec).exec,
		&streamq.TSCreate{}:   newStreamCreate(s.qExec, s.obs, s.rngSvc).exec,
	})
}

func (s *Service) NewTSCreate() *StreamCreate {
	return newStreamCreate(s.qExec, s.obs, s.rngSvc)
}

func (s *Service) NewTSRetrieve() *StreamRetrieve {
	return newStreamRetrieve(s.qExec)
}
