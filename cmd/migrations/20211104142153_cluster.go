package migrations

import (
	"context"
	"github.com/arya-analytics/aryacore/cluster"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().Model((*cluster.Node)(nil)).Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewDropTable().Model((*cluster.Node)(nil)).Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	})
}
