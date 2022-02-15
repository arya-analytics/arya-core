package cluster

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Service interface {
	CanHandle(q *Query) bool
	Exec(ctx context.Context, q *Query) error
}

type QueryVariant int

const (
	QueryVariantCreate = iota
)

type QueryOpt func() interface{}

type QueryOpts map[string]QueryOpt

func (qp QueryOpts) Retrieve(key string) (QueryOpt, bool) {
	q, ok := qp[key]
	return q, ok
}

type Query struct {
	Variant QueryVariant
	Model   *model.Reflect
	Params  QueryOpts
}

type Cluster struct {
	storage  storage.Storage
	services []Service
}

func New(storage storage.Storage) *Cluster {
	return &Cluster{storage: storage}
}
