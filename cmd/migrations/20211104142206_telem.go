package migrations

import (
	"context"
	telem2 "github.com/arya-analytics/aryacore/pkg/telem"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewCreateTable().
			Model((*telem2.ChannelConfig)(nil)).
			ForeignKey(`("node_id") REFERENCES "nodes" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {

			panic(err)
		}
		if _, err := db.NewCreateTable().
			Model((*telem2.RangeReplicaToNode)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewCreateTable().
			Model((*telem2.Range)(nil)).
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
			Model((*telem2.ChannelChunk)(nil)).
			ForeignKey(`("channel_config_id") REFERENCES "channel_configs" ("id") 
						ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		if _, err := db.NewDropTable().
			Model((*telem2.ChannelChunk)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem2.RangeReplicaToNode)(nil)).
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem2.Range)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		if _, err := db.NewDropTable().
			Model((*telem2.ChannelConfig)(nil)).
			IfExists().
			Exec(ctx); err != nil {
			panic(err)
		}
		return nil
	})
}
