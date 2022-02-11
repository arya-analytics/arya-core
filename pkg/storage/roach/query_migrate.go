package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type queryMigrate struct {
	queryBase
	db         *bun.DB
	migrations *bunMigrate.Migrations
	driver     Driver
}

func newMigrate(db *bun.DB, driver Driver) *queryMigrate {
	q := &queryMigrate{
		db:         db,
		migrations: bunMigrate.NewMigrations(),
		driver:     driver,
	}
	q.baseInit(db)
	bindMigrations(q.migrations, q.driver)
	return q
}

func (q *queryMigrate) bunMigrator() *bunMigrate.Migrator {
	return bunMigrate.NewMigrator(q.db, q.migrations)
}

func (q *queryMigrate) init(ctx context.Context) {
	q.catcher.Exec(func() error { return q.bunMigrator().Init(ctx) })
}

func (q *queryMigrate) Exec(ctx context.Context) error {
	q.init(ctx)
	var group *bunMigrate.MigrationGroup
	q.catcher.Exec(func() (err error) {
		group, err = q.bunMigrator().Migrate(ctx)
		return err
	})
	q.catcher.Exec(func() error {
		if group.ID == 0 {
			log.Info("No new migrations to run.")
		}
		log.Infof("Migrated to group %s \n", group)
		return nil
	})
	return q.baseErr()
}

func (q *queryMigrate) Verify(ctx context.Context) error {
	q.baseExec(func() error {
		_, err := q.db.NewSelect().Model((*ChannelConfig)(nil)).Count(ctx)
		return err
	})
	return q.baseErr()
}
