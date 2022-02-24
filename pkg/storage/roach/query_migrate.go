package roach

import (
	"context"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type queryMigrate struct {
	queryBase
	db     *bun.DB
	bunQ   *bunMigrate.Migrations
	driver Driver
}

func newMigrate(db *bun.DB, driver Driver) *queryMigrate {
	q := &queryMigrate{
		db:     db,
		bunQ:   bunMigrate.NewMigrations(),
		driver: driver,
	}
	q.baseInit(db)
	bindMigrations(q.bunQ, q.driver)
	return q
}

func (q *queryMigrate) bunMigrator() *bunMigrate.Migrator {
	return bunMigrate.NewMigrator(q.db, q.bunQ)
}

func (q *queryMigrate) init(ctx context.Context) {
	q.catcher.Exec(func() error { return q.bunMigrator().Init(ctx) })
}

func (q *queryMigrate) Exec(ctx context.Context) error {
	q.init(ctx)
	q.catcher.Exec(func() error {
		_, err := q.bunMigrator().Migrate(ctx)
		return err
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
