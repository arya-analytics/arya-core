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
	// CRDB Internal Schema.
	crdbSchema         = "crdb_internal"
	crdbGossipNodes    = crdbSchema + ".gossip_nodes"
	crdbGossipLiveness = crdbSchema + ".gossip_liveness"
)

var (
	/* If we're using DriverPG, we expect two CRDB internal tables that provide
	information on Node identity and status. This logic creates a view so that this
	info can be accessed by the ORM. */
	nodesViewSQL = fmt.Sprintf(`CREATE VIEW %s AS SELECT n.id,
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

		// |||| NODE ||||

		c.execMigration(db.NewCreateTable().Model((*Node)(nil)).Exec)
		c.Exec(func() error {
			_, err := db.Exec(nodesViewSQL)
			return err
		})

		// |||| RANGE ||||
		c.execMigration(db.NewCreateTable().
			Model((*Range)(nil)).
			Exec,
		)
		c.execMigration(db.NewCreateTable().
			Model((*RangeReplica)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.execMigration(db.NewCreateTable().
			Model((*RangeLease)(nil)).
			ForeignKey(`("range_replica_id") REFERENCES "range_replicas" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.Exec(func() error {
			_, err := db.Exec(`ALTER TABLE "ranges" ADD CONSTRAINT fk_range_lease_id_ref_range_leases FOREIGN KEY ("range_lease_id") REFERENCES "range_leases" ("id") ON DELETE CASCADE`)
			return err
		})

		// |||| CHANNEL ||||

		c.execMigration(db.NewCreateTable().
			Model((*ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.execMigration(db.NewCreateTable().
			Model((*ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.execMigration(db.NewCreateTable().
			Model((*ChannelChunkReplica)(nil)).
			ForeignKey(`("channel_chunk_id") REFERENCES channel_chunks ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_replica_id") REFERENCES range_replicas ("id") ON DELETE CASCADE`).
			Exec,
		)

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
