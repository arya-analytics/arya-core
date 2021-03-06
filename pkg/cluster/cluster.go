package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service interface {
	CanHandle(q *query.Pack) bool
	query.AssembleExec
}

type Cluster interface {
	query.Assemble
	BindService(s Service)
}

type cluster struct {
	query.Assemble
	svc ServiceChain
}

func New() Cluster {
	c := &cluster{}
	c.Assemble = query.NewAssemble(c.Exec)
	return c
}

func (c *cluster) BindService(s Service) {
	c.svc = append(c.svc, s)
}

func (c *cluster) Exec(ctx context.Context, p *query.Pack) error {
	return c.svc.Exec(ctx, p)
}
