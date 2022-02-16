package cluster

import (
	"context"
)

type Service interface {
	CanHandle(q *QueryRequest) bool
	Exec(ctx context.Context, q *QueryRequest) error
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
