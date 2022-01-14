package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
	"reflect"
)

type migrate struct {
	db         *bun.DB
	migrations *bunMigrate.Migrations
	driver     Driver
}

func newMigrate(db *bun.DB, driver Driver) *migrate {
	m := &migrate{
		db:         db,
		migrations: bunMigrate.NewMigrations(),
		driver:     driver,
	}
	bindMigrations(m.migrations, m.driver)
	return m
}

func (m *migrate) bunMigrator() *bunMigrate.Migrator {
	return bunMigrate.NewMigrator(m.db, m.migrations)
}

func (m *migrate) init(ctx context.Context) error {
	return m.bunMigrator().Init(ctx)
}

func (m *migrate) Exec(ctx context.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}
	group, err := m.bunMigrator().Migrate(ctx)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		log.Info("No new migrations to run.")
	}
	log.Infof("Migrated to group %s \n", group)
	return nil
}

func (m *migrate) Verify(ctx context.Context) (err error) {
	for _, rm := range models() {
		_, err = m.db.NewSelect().Model(reflect.New(rm).Interface()).Count(ctx)
	}
	log.Warn(err)
	return err
}
