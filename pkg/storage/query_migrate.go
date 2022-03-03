package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type QueryMigrate struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newMigrate(s *storage) *QueryMigrate {
	q := &QueryMigrate{}
	q.baseInit(s, q)
	return q
}

/// |||| INTERFACE ||||

func (q *QueryMigrate) Exec(ctx context.Context) error {
	c := errutil.NewCatchWCtx(ctx)
	c.Exec(q.storage.cfg.EngineMD.NewMigrate().Exec)
	c.Exec(q.storage.cfg.EngineObject.NewMigrate().Exec)
	return q.baseErr()
}
