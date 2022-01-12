package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type migrator struct {
	db         *bun.DB
	migrations *migrate.Migrations
}

func newMigrator(db *bun.DB) *migrator {
	m := &migrator{
		db:         db,
		migrations: migrate.NewMigrations(),
	}
	bindMigrations(m.migrations)
	return m
}

func (m *migrator) migrator() *migrate.Migrator {
	return migrate.NewMigrator(m.db, m.migrations)
}

func (m *migrator) init(ctx context.Context) error {
	return m.migrator().Init(ctx)
}

func (m *migrator) migrate(ctx context.Context) error {
	group, err := m.migrator().Migrate(ctx)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		log.Info("No new migrations to run.")
	}
	log.Infof("Migrated to group %s \n", group)
	return nil
}

func (m *migrator) verify(ctx context.Context) (err error) {
	for _, rm := range models {
		_, err = m.db.NewSelect().Model(rm).Count(ctx)
	}
	return err
}
