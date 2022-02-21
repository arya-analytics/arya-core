package storage

import "context"

type QueryMigrate struct {
	queryBase
}

// |||| CONSTRUCTOR ||||

func newMigrate(s Storage) *QueryMigrate {
	q := &QueryMigrate{}
	q.baseInit(s, s.config().Hooks.BeforeMigrate)
	return q
}

/// |||| INTERFACE ||||

func (q *QueryMigrate) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		return q.mdQuery().Exec(ctx)
	})
	q.baseExec(func() error {
		return q.objQuery().Exec(ctx)
	})
	return q.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (q *QueryMigrate) mdQuery() QueryMDMigrate {
	if q.baseMDQuery() == nil {
		q.baseSetMDQuery(q.baseMDEngine().NewMigrate(q.baseMDAdapter()))
	}
	return q.baseMDQuery().(QueryMDMigrate)
}

// || OBJECT ||

func (q *QueryMigrate) objQuery() QueryObjectMigrate {
	if q.baseObjQuery() == nil {
		q.baseSetObjQuery(q.baseObjEngine().NewMigrate(q.baseObjAdapter()))
	}
	return q.baseObjQuery().(QueryObjectMigrate)
}
