package migrations

import (
	"context"
	"github.com/arya-analytics/aryacore/telem"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().
			Model((*telem.ChannelConfig)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewCreateTable().
			Model((*telem.RangeReplicaToNode)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewCreateTable().
			Model((*telem.Range)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.Exec(`ALTER TABLE range_replica_to_nodes ADD CONSTRAINT 
									range_replica_to_nodes_fk_node_id FOREIGN KEY 
									("node_id") REFERENCES "nodes" ("id") ON DELETE 
									CASCADE`); err != nil {
			panic(err)
		}
		if _, err := db.Exec(`ALTER TABLE range_replica_to_nodes ADD CONSTRAINT 
									range_replica_to_nodes FOREIGN KEY 
									("range_id") REFERENCES "ranges" ("id") ON DELETE 
									CASCADE`); err != nil {
			panic(err)
		}
		if _, err := db.NewCreateTable().
			Model((*telem.ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") 
						ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewDropTable().
			Model((*telem.ChannelChunk)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem.RangeReplicaToNode)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem.Range)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem.ChannelConfig)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	})
}
