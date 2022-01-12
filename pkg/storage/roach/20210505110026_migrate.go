package roach

import (
	"context"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func migrateUp(ctx context.Context, db *bun.DB) error {
	if _, err := db.NewCreateTable().Model((*ChannelConfig)(nil)).Exec(
		ctx); err != nil {
		panic(err)
	}
	return nil
}

func migrateDown(ctx context.Context, db *bun.DB) error {
	return nil
}

func bindMigrations(m *migrate.Migrations) {
	m.MustRegister(migrateUp, migrateDown)
}
