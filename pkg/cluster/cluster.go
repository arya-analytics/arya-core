package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type Service interface {
	CanHandle(q *internal.QueryRequest) bool
	Exec(ctx context.Context, q *internal.QueryRequest) error
}

type Cluster interface {
	NewCreate() *QueryCreate
	NewRetrieve() *QueryRetrieve
	NewUpdate() *QueryUpdate
	NewDelete() *QueryDelete
}

type cluster struct {
	svc ServiceChain
}

func New(svc ServiceChain) Cluster {
	return &cluster{svc}
}

func (c *cluster) NewCreate() *QueryCreate {
	return newCreate(c.svc)
}

func (c *cluster) NewRetrieve() *QueryRetrieve {
	return newRetrieve(c.svc)
}

func (c *cluster) NewUpdate() *QueryUpdate {
	return newUpdate(c.svc)
}

func (c *cluster) NewDelete() *QueryDelete {
	return newDelete(c.svc)
}
