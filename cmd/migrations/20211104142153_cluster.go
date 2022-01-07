package migrations

import (
	"context"
	cluster2 "github.com/arya-analytics/aryacore/pkg/cluster"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().Model((*cluster2.Node)(nil)).Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewDropTable().Model((*cluster2.Node)(nil)).Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	})
}
