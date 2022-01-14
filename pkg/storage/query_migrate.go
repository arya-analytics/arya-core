package storage

import "context"

type migrateQuery struct {
	baseQuery
	_mdQuery MigrateQuery
}

func newMigrate(s *Storage) *migrateQuery {
	m := &migrateQuery{}
	m.init(s)
	return m
}

func (m *migrateQuery) mdQuery() MigrateQuery {
	if m._mdQuery == nil {
		m._mdQuery = m.mdEngine.NewMigrate(m.storage.adapter(EngineRoleMD))
	}
	return m._mdQuery
}

func (m *migrateQuery) Exec(ctx context.Context) error {
	return m.mdQuery().Exec(ctx)
}

