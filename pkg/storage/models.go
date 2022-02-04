package storage

import (
	"github.com/google/uuid"
	"time"
)

type Node struct {
	ID              int `model:"role:pk"`
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
	ID                uuid.UUID `model:"role:pk"`
	LeaseHolderNode   *Node
	LeaseHolderNodeID int
	ReplicaNodes      []*Node
}

type RangeReplicaToNode struct {
	ID      uuid.UUID `model:"role:pk"`
	Range   *Range
	RangeID uuid.UUID
	Node    *Node
	NodeID  int
}

type ChannelConfig struct {
	ID        uuid.UUID `model:"role:pk"`
	Name      string
	Node      *Node
	NodeID    int
	DataRate  float64
	Retention time.Duration
}

type ChannelChunk struct {
	ID              uuid.UUID `model:"role:pk"`
	Range           *Range
	RangeID         uuid.UUID
	ChannelConfig   *ChannelConfig
	ChannelConfigID uuid.UUID
	Data            Object `storage:"re:object"`
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `model:"role:series"`
	ChannelConfigID uuid.UUID
	Value           float64 `storage:"role:cache"`
	Timestamp       int64   `storage:"role:cache"`
}
