package roach

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

func migrateUpFunc(d Driver) bunMigrate.MigrationFunc {
	return func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().Model((*Node)(nil)).
			Exec(ctx); err != nil {
			log.Fatalln(err)
		}
		if d == DriverPG {
			if _, err := db.Exec(`CREATE VIEW nodes_w_gossip AS SELECT n.id, 
									gn.address, gn.is_live, gn.started_at, gv.epoch, 
									gv.draining, gv.decommissioning, gv.membership, 
									gv.updated_at FROM nodes n 
									JOIN crdb_internal.gossip_nodes gn ON n.id =
									gn.node_id LEFT JOIN crdb_internal.
									gossip_liveness gv ON gv.node_id=n.id`); err != nil {
				log.Fatalln(err)
			}
		} else if d == DriverSQLite {
			if _, err := db.Exec(`CREATE VIEW nodes_w_gossip AS SELECT n.id
									FROM nodes n`); err != nil {
				log.Fatalln(err)
			}

		}
		if _, err := db.NewCreateTable().
			Model((*ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {

			log.Fatalln(err)
		}
		if _, err := db.NewCreateTable().
			Model((*RangeReplicaToNode)(nil)).
			Exec(ctx); err != nil {
			log.Fatalln(err)
		}
		if _, err := db.NewCreateTable().
			Model((*Range)(nil)).
			ForeignKey(`("lease_holder_node_id") REFERENCES "nodes" (
						"id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			log.Fatalln(err)
		}
		//if _, err := db.Exec(`ALTER TABLE range_replica_to_nodes ADD FOREIGN KEY
		//							("node_id") REFERENCES "nodes" ("id") ON DELETE
		//							CASCADE`); err != nil {
		//	log.Fatalln(err)
		//}
		//if _, err := db.Exec(`ALTER TABLE range_replica_to_nodes ADD FOREIGN KEY
		//							("range_id") REFERENCES "ranges" ("id") ON DELETE
		//							CASCADE`); err != nil {
		//	log.Fatalln(err)
		//}
		if _, err := db.NewCreateTable().
			Model((*ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") 
						ON DELETE CASCADE`).
			Exec(ctx); err != nil {
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
