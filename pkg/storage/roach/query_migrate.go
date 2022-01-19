package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type migrateQuery struct {
	baseQuery
	db         *bun.DB
	migrations *bunMigrate.Migrations
	driver     Driver
}

func newMigrate(db *bun.DB, driver Driver) *migrateQuery {
	m := &migrateQuery{
		db:         db,
		migrations: bunMigrate.NewMigrations(),
		driver:     driver,
	}
	bindMigrations(m.migrations, m.driver)
	return m
}

func (m *migrateQuery) bunMigrator() *bunMigrate.Migrator {
	return bunMigrate.NewMigrator(m.db, m.migrations)
}

func (m *migrateQuery) init(ctx context.Context) error {
	return m.bunMigrator().Init(ctx)
}

func (m *migrateQuery) Exec(ctx context.Context) error {
	if err := m.init(ctx); err != nil {
		return m.baseHandleExecErr(err)
	}
	group, err := m.bunMigrator().Migrate(ctx)
	if err != nil {
		return m.baseHandleExecErr(err)
	}
	if group.ID == 0 {
		log.Info("No new migrations to run.")
	}
	log.Infof("Migrated to group %s \n", group)
	return nil
}

func (m *migrateQuery) Verify(ctx context.Context) (err error) {
	_, err = m.db.NewSelect().Model((*ChannelConfig)(nil)).Count(ctx)
	return m.baseHandleExecErr(err)
}
