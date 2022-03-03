package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service interface {
	CanHandle(q *query.Pack) bool
	Exec(ctx context.Context, p *query.Pack) error
}

type Cluster interface {
	query.Assemble
	BindService(s Service)
	Exec(ctx context.Context, p *query.Pack) error
}

type cluster struct {
	query.AssembleBase
	svc ServiceChain
}

func New() Cluster {
	c := &cluster{}
	c.AssembleBase = query.NewAssemble(c.Exec)
	return c
}

func (c *cluster) BindService(s Service) {
	c.svc = append(c.svc, s)
}

func (c *cluster) Exec(ctx context.Context, p *query.Pack) error {
	return c.svc.Exec(ctx, p)
}
