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
	if err := m.mdQuery().Exec(ctx); err != nil {
		return err
	}
	if err := m.objQuery().Exec(ctx); err != nil {
		return err
	}
	return nil
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
