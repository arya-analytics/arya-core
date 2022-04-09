package roach

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

const (
	nodesGossip = "nodes_gossip"
	// CRDB Internal Schema.
	crdbSchema         = "crdb_internal"
	crdbGossipNodes    = crdbSchema + ".gossip_nodes"
	crdbGossipLiveness = crdbSchema + ".gossip_liveness"
	crdNodeRuntimeInfo = crdbSchema + ".node_runtime_info"
)

var (
	/* If we're using DriverPG, we expect two CRDB internal tables that provide
	information on Node identity and status. This logic creates a view so that this
	info can be accessed by the ORM. */
	nodesViewSQL = fmt.Sprintf(`CREATE VIEW %s AS SELECT DISTINCT n.id,
									gn.address,
									gn.is_live, 
									gn.started_at, 
									gv.epoch, 
									gv.draining, 
									gv.decommissioning, 
									gv.membership, 
									gv.updated_at,
									n.rpc_port,
									n.id = nri.node_id is_host
									FROM nodes n 
									JOIN %s gn ON n.id = gn.node_id 
									LEFT JOIN %s gv ON gv.node_id=n.id
									LEFT JOIN %s nri ON nri.node_id=n.id`,
		nodesGossip,
		crdbGossipNodes,
		crdbGossipLiveness,
		crdNodeRuntimeInfo,
	)
)

// |||| CATCHER ||||

type migrateCatcher struct {
	*errutil.CatchContext
}

type migrationExecFunc func(ctx context.Context, dest ...interface{}) (sql.Result, error)

func (m *migrateCatcher) Exec(execFunc migrationExecFunc) {
	m.CatchContext.Exec(func(ctx context.Context) error {
		_, err := execFunc(ctx)
		return err
	})
}

// |||| MIGRATE UP ||||

func migrateUpFunc() bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {

		c := &migrateCatcher{CatchContext: errutil.NewCatchContext(ctx)}

		// |||| NODE ||||

		c.Exec(db.NewCreateTable().Model((*models.Node)(nil)).Exec)
		c.CatchSimple.Exec(func() error {
			_, err := db.Exec(nodesViewSQL)
			return err
		})

		// |||| RANGE ||||

		c.Exec(db.NewCreateTable().
			Model((*models.Range)(nil)).
			Exec,
		)
		c.Exec(db.NewCreateIndex().
			Model((*models.Range)(nil)).
			Column("id").
			Where("status > 1").
			Exec,
		)

		c.Exec(db.NewCreateTable().
			Model((*models.RangeReplica)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.Exec(db.NewCreateTable().
			Model((*models.RangeLease)(nil)).
			ForeignKey(`("range_replica_id") REFERENCES "range_replicas" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec,
		)

		// |||| CHANNEL ||||

		c.Exec(db.NewCreateTable().
			Model((*models.ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.Exec(db.NewCreateTable().
			Model((*models.ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_id") REFERENCES "ranges" ("id") ON DELETE CASCADE`).
			Exec,
		)
		c.Exec(db.NewCreateIndex().
			Model((*models.ChannelChunk)(nil)).
			Column("id", "start_ts").
			Exec,
		)
		c.Exec(db.NewCreateTable().
			Model((*models.ChannelChunkReplica)(nil)).
			ForeignKey(`("channel_chunk_id") REFERENCES channel_chunks ("id") ON DELETE CASCADE`).
			ForeignKey(`("range_replica_id") REFERENCES range_replicas ("id") ON DELETE CASCADE`).
			Exec,
		)

		// |||| AUTH ||||

		db.RegisterModel((*models.UserToGroup)(nil))

		c.Exec(db.NewCreateTable().Model((*models.User)(nil)).Exec)
		c.Exec(db.NewCreateTable().Model((*models.Group)(nil)).Exec)
		c.Exec(db.NewCreateTable().
			Model((*models.UserToGroup)(nil)).
			ForeignKey(`("group_id") REFERENCES groups ("id") ON DELETE CASCADE`).
			ForeignKey(`("user_id") REFERENCES users ("id") ON DELETE CASCADE`).
			Exec,
		)

		return c.Error()
	}
}

// |||| MIGRATE DOWN ||||

func migrateDownFunc() bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		return nil
	}
}

// |||| MIGRATION BINDING ||||

func bindMigrations(m *bunMigrate.Migrations) {
	m.MustRegister(migrateUpFunc(), migrateDownFunc())
}
