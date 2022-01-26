package roach

import (
	"context"
	"database/sql"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type migrateCatcher struct {
	ctx context.Context
	*errutil.Catcher
}

type MigrationExecFunc func(ctx context.Context, dest ...interface{}) (sql.Result, error)

func (m *migrateCatcher) ExecMigration(execFunc MigrationExecFunc) {
	m.Exec(func() error {
		_, err := execFunc(m.ctx)
		return err
	})
}

func migrateUpFunc(d Driver) bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		c := &migrateCatcher{Catcher: &errutil.Catcher{}, ctx: ctx}
		db.RegisterModel((*RangeReplicaToNode)(nil))
		c.ExecMigration(db.NewCreateTable().Model((*Node)(nil)).Exec)
		if d == DriverPG {
			c.Exec(func() error {
				_, err := db.Exec(
					`CREATE VIEW nodes_w_gossip AS SELECT n.id,
									gn.address, gn.is_live, gn.started_at, gv.epoch, 
									gv.draining, gv.decommissioning, gv.membership, 
									gv.updated_at FROM nodes n 
									JOIN crdb_internal.gossip_nodes gn ON n.id =
									gn.node_id LEFT JOIN crdb_internal.
									gossip_liveness gv ON gv.node_id=n.id`)
				return err
			})
		} else if d == DriverSQLite {
			c.Exec(func() error {
				_, err := db.Exec(`CREATE VIEW nodes_w_gossip AS SELECT n.id
									FROM nodes n`)
				return err
			})

		}
		c.ExecMigration(db.NewCreateTable().
			Model((*ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec)

		c.ExecMigration(db.NewCreateTable().
			Model((*Range)(nil)).
			ForeignKey(`("lease_holder_node_id") REFERENCES "nodes" (
						"id") ON DELETE CASCADE`).
			Exec)
		c.ExecMigration(db.NewCreateTable().
			Model((*RangeReplicaToNode)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec)

		c.ExecMigration(db.NewCreateTable().
			Model((*ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") 
						ON DELETE CASCADE`).
			Exec)
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
