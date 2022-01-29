package storage

import (
	"github.com/google/uuid"
	"time"
)

var _catalog = ModelCatalog{
	&Node{},
	&Range{},
	&RangeReplicaToNode{},
	&ChannelConfig{},
	&ChannelChunk{},
	&ChannelSample{},
}

func catalog() ModelCatalog {
	return _catalog
}

type Node struct {
	ID              int
	Address         string
	StartedAt       time.Time
	IsLive          bool
	Epoch           int
	Expiration      string
	Draining        bool
	Decommissioning bool
	Membership      string
	UpdatedAt       time.Time
}

type Range struct {
	ID uuid.UUID
	// LeaseHolderNode
	LeaseHolderNode   *Node
	LeaseHolderNodeID int
	// ReplicaNodes
	ReplicaNodes []*Node
}

type RangeReplicaToNode struct {
	ID uuid.UUID
	// Range
	Range   *Range
	RangeID uuid.UUID
	// Node
	Node   *Node
	NodeID int
}

type ChannelConfig struct {
	ID   uuid.UUID
	Name string
	// Node
	Node   *Node
	NodeID int
	// Data
	DataRate  float64
	Retention time.Duration
}

type ChannelChunk struct {
	ID uuid.UUID
	// Range
	Range   *Range
	RangeID uuid.UUID
	// ChannelConfig
	ChannelConfig   *ChannelConfig
	ChannelConfigID uuid.UUID
	// Data
	Data Object
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `storage:"role:index"`
	ChannelConfigID uuid.UUID
	Value           float64
	Timestamp       int64
}
