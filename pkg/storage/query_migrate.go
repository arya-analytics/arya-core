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
	return m.mdQuery().Exec(ctx)
}

// |||| QUERY BINDING ||||

func (m *migrateQuery) mdQuery() MDMigrateQuery {
	if m.baseMDQuery() == nil {
		m.baseSetMDQuery(m.mdEngine.NewMigrate(m.storage.adapter(EngineRoleMD)))
	}
	return m.baseMDQuery().(MDMigrateQuery)
}
