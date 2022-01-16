package storage

import (
	"github.com/google/uuid"
	"time"
)

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
	ID                int
	LeaseHolderNodeID int
	LeaseHolderNode   *Node
	ReplicaNodes      []*Node
}
type RangeReplicaToNode struct {
	ID      uuid.UUID
	RangeID int
	Range   *Range
	NodeID  int
	Node    *Node
}

type ChannelConfig struct {
	ID     int
	Name   string
	NodeID int
	Node   *Node
}

type ChannelChunk struct {
	ID              int
	RangeID         int
	Range           *Range
	ChannelConfigID int
	ChannelConfig   *Range
}
