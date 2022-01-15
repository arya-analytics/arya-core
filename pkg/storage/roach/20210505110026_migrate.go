package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

func migrateUpFunc(d Driver) bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().Model((*ChannelConfig)(nil)).Exec(
			ctx); err != nil {
			log.Fatalln(err)
		}
		return nil
	}
}

func migrateDownFunc(d Driver) bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		return nil
	}
}

func bindMigrations(m *bunMigrate.Migrations, d Driver) {
	m.MustRegister(migrateUpFunc(d), migrateDownFunc(d))
}
