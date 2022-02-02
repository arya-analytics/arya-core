package storage

import "context"

type MigrateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newMigrate(s *Storage) *MigrateQuery {
	m := &MigrateQuery{}
	m.baseInit(s)
	return m
}

/// |||| INTERFACE ||||

func (m *MigrateQuery) Exec(ctx context.Context) error {
	m.baseExec(func() error {
		return m.mdQuery().Exec(ctx)
	})
	m.baseExec(func() error {
		return m.objQuery().Exec(ctx)
	})
	return m.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (m *MigrateQuery) mdQuery() MDMigrateQuery {
	if m.baseMDQuery() == nil {
		m.baseSetMDQuery(m.storage.cfg.mdEngine().NewMigrate(m.baseMDAdapter()))
	}
	return m.baseMDQuery().(MDMigrateQuery)
}

// || OBJECT ||

func (m *MigrateQuery) objQuery() ObjectMigrateQuery {
	if m.baseObjQuery() == nil {
		m.baseSetObjQuery(m.storage.cfg.objEngine().NewMigrate(m.baseObjAdapter()))
	}
	return m.baseObjQuery().(ObjectMigrateQuery)
}
