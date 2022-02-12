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
	IsHost          bool
	Epoch           int
	Expiration      string
	Draining        bool
	Decommissioning bool
	Membership      string
	UpdatedAt       time.Time
}

type Range struct {
	ID           uuid.UUID `model:"role:pk,"`
	RangeLeaseID uuid.UUID
}

type RangeLease struct {
	ID             uuid.UUID `model:"role:pk,"`
	Range          *Range
	RangeReplica   *RangeReplica
	RangeReplicaID uuid.UUID
}

type RangeReplica struct {
	ID      uuid.UUID `model:"role:pk"`
	Range   *Range
	RangeID uuid.UUID
	Node    *Node
	NodeID  int
}

// |||| CHANNEL ||||

type ChannelConfig struct {
	ID     uuid.UUID `model:"role:pk,"`
	Name   string
	Node   *Node
	NodeID int
}

type ChannelChunk struct {
	ID              uuid.UUID `model:"role:pk,"`
	Range           *Range
	RangeID         uuid.UUID
	ChannelConfig   *ChannelConfig
	ChannelConfigID uuid.UUID
}

type ChannelChunkReplica struct {
	ID             uuid.UUID `model:"role:pk,"`
	ChannelChunk   *ChannelChunk
	ChannelChunkID uuid.UUID
	RangeReplica   *RangeReplica
	RangeReplicaID uuid.UUID
	Telem          Object `storage:"re:object,role:bulkTelem"`
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `model:"role:series"`
	ChannelConfigID uuid.UUID
	Value           float64 `storage:"role:cache"`
	Timestamp       int64   `storage:"role:cache"`
}
