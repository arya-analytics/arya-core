package storage

import "context"

type migrateQuery struct {
	baseQuery
}

func newMigrate(s *Storage) *migrateQuery {
	m := &migrateQuery{}
	m.baseInit(s)
	return m
}

func (m *migrateQuery) mdQuery() MDMigrateQuery {
	if m.baseMDQuery() == nil {
		m.baseSetMDQuery(m.mdEngine.NewMigrate(m.storage.adapter(EngineRoleMD)))
	}
	return m.baseMDQuery().(MDMigrateQuery)
}

func (m *migrateQuery) Exec(ctx context.Context) error {
	return m.mdQuery().Exec(ctx)
}
