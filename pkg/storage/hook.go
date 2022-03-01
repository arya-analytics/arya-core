package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type QueryHook interface {
	BeforeQuery(ctx context.Context, o *query.Pack) error
	AfterQuery(ctx context.Context, qe *query.Pack) error
}
