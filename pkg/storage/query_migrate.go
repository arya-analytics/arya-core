package storage

import "context"

type migrateQuery struct {
	baseQuery
}

// |||| CONSTRUCTOR ||||

func newMigrate(s *Storage) *migrateQuery {
	m := &migrateQuery{}
	m.baseInit(s)
	return m
}

/// |||| INTERFACE ||||

func (m *migrateQuery) Exec(ctx context.Context) error {
	m.catcher.Exec(func() error {
		return m.mdQuery().Exec(ctx)
	})
	m.catcher.Exec(func() error {
		return m.objQuery().Exec(ctx)
	})
	return m.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (m *migrateQuery) mdQuery() MDMigrateQuery {
	if m.baseMDQuery() == nil {
		m.baseSetMDQuery(m.mdEngine.NewMigrate(m.baseMDAdapter()))
	}
	return m.baseMDQuery().(MDMigrateQuery)
}

// || OBJECT ||

func (m *migrateQuery) objQuery() ObjectMigrateQuery {
	if m.baseObjQuery() == nil {
		m.baseSetObjQuery(m.objEngine.NewMigrate(m.baseObjAdapter()))
	}
	return m.baseObjQuery().(ObjectMigrateQuery)
}
