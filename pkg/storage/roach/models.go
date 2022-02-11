package roach

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{
		&Node{},
		&Range{},
		&RangeReplica{},
		&RangeLease{},
		&ChannelConfig{},
		&ChannelChunk{},
		&GossipNode{},
		&GossipLiveness{},
		&ChannelChunkReplica{},
	}
}

// |||| DEFINITIONS ||||

// |||| NODE ||||

type Node struct {
	// Select key MUST match nodesGossip table in migrations file.
	bun.BaseModel `bun:"select:nodes_gossip,table:nodes"`
	ID            int `bun:",pk" model:"role:pk,"`
	GossipNode
	GossipLiveness
}

// |||| RANGE ||||

type Range struct {
	ID           uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	RangeLeaseID uuid.UUID `bun:"type:UUID,nullzero"`
}

type RangeLease struct {
	ID             uuid.UUID     `bun:"type:UUID,pk" model:"role:pk,"`
	Range          *Range        `bun:"rel:has-one,join:id=range_lease_id"`
	RangeReplica   *RangeReplica `bun:"rel:belongs-to,join:range_replica_id=id"`
	RangeReplicaID uuid.UUID     `bun:"type:UUID,"`
}

type RangeReplica struct {
	ID      uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	Range   *Range    `bun:"rel:belongs-to,join:range_id=id"`
	RangeID uuid.UUID `bun:"type:UUID,"`
	Node    *Node     `bun:"rel:belongs-to,join:node_id=id"`
	NodeID  int
}

// |||| CHANNEL ||||

type ChannelConfig struct {
	ID     uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	Name   string
	Node   *Node `bun:"rel:belongs-to,join:node_id=id,"`
	NodeID int
}

type ChannelChunk struct {
	ID              uuid.UUID      `bun:"type:UUID,pk" model:"role:pk,"`
	Range           *Range         `bun:"rel:belongs-to,join:range_id=id,"`
	RangeID         uuid.UUID      `bun:"type:UUID,"`
	ChannelConfig   *ChannelConfig `bun:"rel:belongs-to,join:channel_config_id=id,"`
	ChannelConfigID uuid.UUID      `bun:"type:UUID,"`
}

type ChannelChunkReplica struct {
	ID             uuid.UUID     `bun:"type:UUID,pk" model:"role:pk,"`
	ChannelChunk   *ChannelChunk `bun:"rel:belongs-to,join:channel_chunk_id=id,"'`
	ChannelChunkID uuid.UUID     `bun:"type:UUID,"`
	RangeReplica   *RangeReplica `bun:"rel:belongs-to,join:range_replica_id=id,"`
	RangeReplicaID uuid.UUID     `bun:"type:UUID,"`
}

// || ROACH INTERNAL MODELS ||

// GossipNode lives in crdb's internal schema and tracks the nodes in the roach cluster.
type GossipNode struct {
	Address   string    `bun:"type:text,scanonly"`
	StartedAt time.Time `bun:"type:timestamp,scanonly"`
	IsLive    bool      `bun:"type:boolean,scanonly"`
}

// GossipLiveness lives in crdb's internal schema and tracks the health of nodes in
// the roach cluster.
type GossipLiveness struct {
	Epoch           int       `bun:"type:bigint,scanonly"`
	Expiration      string    `bun:"type:text,scanonly"`
	Draining        bool      `bun:"type:boolean,scanonly"`
	Decommissioning bool      `bun:"type:boolean,scanonly"`
	Membership      string    `bun:"type:text,scanonly"`
	UpdatedAt       time.Time `bun:"type:timestamp,scanonly"`
}
