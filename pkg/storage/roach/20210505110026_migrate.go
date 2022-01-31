package roach

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

const (
	nodesGossip = "nodes_gossip"
	// CRDB Internal Schema
	crdbSchema         = "crdb_internal"
	crdbGossipNodes    = crdbSchema + ".gossip_nodes"
	crdbGossipLiveness = crdbSchema + ".gossip_liveness"
)

var (
	/* If we're using DriverPG, we expect two CRDB internal tables that provide
	information on Node identity and status. This logic creates a view so that this
	info can be accessed by the ORM. */
	driverPGNodesViewSQL = fmt.Sprintf(`CREATE VIEW %s AS SELECT n.id,
									gn.address, 
									gn.is_live, 
									gn.started_at, 
									gv.epoch, 
									gv.draining, 
									gv.decommissioning, 
									gv.membership, 
									gv.updated_at 
									FROM nodes n 
									JOIN %s gn ON n.id =
									gn.node_id LEFT JOIN %s gv ON gv.node_id=n.id`,
		nodesGossip,
		crdbGossipNodes,
		crdbGossipLiveness)
	/* If we're using DriverSQLite, CRDB internal tables aren't available,
	so we just only map the view to the Node table,
	that way we don't need to change any ORM logic. */
	driverSQLiteNodesViewSQL = fmt.Sprintf(`CREATE VIEW %s AS SELECT n.id
									FROM nodes n`, nodesGossip)
)

// |||| CATCHER ||||

type migrateCatcher struct {
	*errutil.Catcher
	ctx context.Context
}

type migrationExecFunc func(ctx context.Context, dest ...interface{}) (sql.Result, error)

func (m *migrateCatcher) execMigration(execFunc migrationExecFunc) {
	m.Exec(func() error {
		_, err := execFunc(m.ctx)
		return err
	})
}

// |||| MIGRATE UP ||||

func migrateUpFunc(d Driver) migrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		c := &migrateCatcher{Catcher: &errutil.Catcher{}, ctx: ctx}
		// Binds the many-to-many relationship the bun ORM,
		// so we can properly run queries against it.
		db.RegisterModel((*rangeReplicaToNode)(nil))

		c.execMigration(db.NewCreateTable().Model((*Node)(nil)).Exec)

		if d == DriverPG {
			c.Exec(func() error {
				_, err := db.Exec(driverPGNodesViewSQL)
				return err

			})
		} else if d == DriverSQLite {

			c.Exec(func() error {
				_, err := db.Exec(driverSQLiteNodesViewSQL)
				return err
			})
		}

		c.execMigration(db.NewCreateTable().
			Model((*ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec)

		c.execMigration(db.NewCreateTable().
			Model((*Range)(nil)).
			ForeignKey(`("lease_holder_node_id") REFERENCES "nodes" (
						"id") ON DELETE CASCADE`).
			Exec)
		c.execMigration(db.NewCreateTable().
			Model((*rangeReplicaToNode)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec)

		c.execMigration(db.NewCreateTable().
			Model((*ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") 
						ON DELETE CASCADE`).
			Exec)
		return c.Error()
	}
}

// |||| MIGRATE DOWN ||||

func migrateDownFunc(d Driver) migrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		return nil
	}
}

// |||| MIGRATION BINDING ||||

func bindMigrations(m *migrate.Migrations, d Driver) {
	m.MustRegister(migrateUpFunc(d), migrateDownFunc(d))
}
