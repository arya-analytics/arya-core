package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type QueryEvent struct {
	Query Query
	Model *model.Reflect
}

type QueryHook interface {
	BeforeQuery(ctx context.Context, qe *QueryEvent) error
	AfterQuery(ctx context.Context, qe *QueryEvent) error
}
