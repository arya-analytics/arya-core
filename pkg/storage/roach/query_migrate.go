package roach

import (
	"context"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type migrateExec struct {
	db     *bun.DB
	bunQ   *bunMigrate.Migrations
	driver Driver
}

func newMigrate(db *bun.DB, driver Driver) *migrateExec {
	q := &migrateExec{
		db:     db,
		bunQ:   bunMigrate.NewMigrations(),
		driver: driver,
	}
	bindMigrations(q.bunQ, q.driver)
	return q
}

func (q *migrateExec) bunMigrator() *bunMigrate.Migrator {
	return bunMigrate.NewMigrator(q.db, q.bunQ)
}

func (q *migrateExec) init(ctx context.Context) error {
	return q.bunMigrator().Init(ctx)
}

func (q *migrateExec) Exec(ctx context.Context) error {
	if err := q.init(ctx); err != nil {
		return q.handleErr(err)
	}
	_, err := q.bunMigrator().Migrate(ctx)
	return q.handleErr(err)
}

func (q *migrateExec) Verify(ctx context.Context) error {
	_, err := q.db.NewSelect().Model((*ChannelConfig)(nil)).Count(ctx)
	return q.handleErr(err)
}

func (q *migrateExec) handleErr(err error) error {
	return newErrorConvert().Exec(err)
}
