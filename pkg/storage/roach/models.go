package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// |||| CATALOG ||||

var _catalog = storage.ModelCatalog{
	&Node{},
	&Range{},
	&rangeReplicaToNode{},
	&ChannelConfig{},
	&ChannelChunk{},
	&GossipNode{},
	&GossipLiveness{},
}

func catalog() storage.ModelCatalog {
	return _catalog
}

// |||| DEFINITIONS ||||

type Node struct {
	// Select key MUST match nodesGossip table in migrations file.
	bun.BaseModel `bun:"select:nodes_gossip,table:nodes" model:"role:pk"`
	ID            int `bun:",pk"`
	GossipNode
	GossipLiveness
}

type Range struct {
	ID                uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	LeaseHolderNodeID int
	LeaseHolderNode   int
	ReplicaNodes      []*Node `bun:"m2m:range_replica_to_nodes,join:Range=Node"`
}

type rangeReplicaToNode struct {
	ID      uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	RangeID uuid.UUID `bun:"type:UUID,"`
	Range   *Range    `bun:"rel:belongs-to,join:range_id=id"`
	NodeID  int
	Node    *Node `bun:"rel:belongs-to,join:node_id=id"`
}

type ChannelConfig struct {
	ID     uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	Name   string
	NodeID int
	Node   *Node `bun:"rel:belongs-to,join:node_id=id,scanonly"`
}

type ChannelChunk struct {
	ID              uuid.UUID      `bun:"type:UUID,pk" model:"role:pk,"`
	RangeID         uuid.UUID      `bun:"type:UUID,"`
	Range           *Range         `bun:"rel:belongs-to,join:range_id=id"`
	ChannelConfigID uuid.UUID      `bun:"type:UUID,"`
	ChannelConfig   *ChannelConfig `bun:"rel:belongs-to,join:channel_config_id=id"`
}

// || ROACH INTERNAL MODELS ||

// GossipNode lives in crdb's internal schema and tracks the nodes in the roach cluster
type GossipNode struct {
	Address   string    `bun:"type:text,scanonly"`
	StartedAt time.Time `bun:"type:timestamp,scanonly"`
	IsLive    bool      `bun:"type:boolean,scanonly"`
}

// GossipLiveness lives in crdb's internal schema and tracks the health of nodes in
//the roach cluster
type GossipLiveness struct {
	Epoch           int       `bun:"type:bigint,scanonly"`
	Expiration      string    `bun:"type:text,scanonly"`
	Draining        bool      `bun:"type:boolean,scanonly"`
	Decommissioning bool      `bun:"type:boolean,scanonly"`
	Membership      string    `bun:"type:text,scanonly"`
	UpdatedAt       time.Time `bun:"type:timestamp,scanonly"`
}
