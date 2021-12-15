package telem

import (
	"github.com/arya-analytics/aryacore/cluster"
	"github.com/google/uuid"
)

type Range struct {
	ID int32 `bun:",pk"`
	string
	LeaseHolderNodeId uuid.UUID
	LeaseHolderNode   *cluster.Node  `bun:"rel:belongs-to,join:lease_holder_node_id=id"`
	ReplicaNodes      []cluster.Node `bun:"m2m:range_replica_to_nodes,join:Range=Node"`
}

type RangeReplicaToNode struct {
	ID      uuid.UUID `bun:"type:uuid,default:gen_random_uuid()"`
	RangeID int32
	Range   *Range        `bun:"rel:belongs-to,join:range_id=id"`
	NodeID  uuid.UUID     `bun:"type:uuid"`
	Node    *cluster.Node `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelConfig struct {
	ID     int32 `bun:",pk"`
	Name   string
	NodeId uuid.UUID     `bun:"type:uuid"`
	Node   *cluster.Node `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelChunk struct {
	ID              int64 `bun:",pk"`
	RangeID         int32
	Range           *Range `bun:"rel:belongs-to,join:range_id=id"`
	ChannelConfigID int32
	ChannelConfig   *Range `bun:"rel:belongs-to,join:channel_config_id=id"`
}
