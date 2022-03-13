package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

type ChannelConflictPolicy int

//go:generate stringer -type=ChannelConflictPolicy
const (
	ChannelConflictPolicyError ChannelConflictPolicy = iota
	ChannelConflictPolicyDiscard
	ChannelConflictPolicyOverwrite
)

type ChannelStatus int

//go:generate stringer -type=ChannelState
const (
	ChannelStatusInactive ChannelStatus = iota + 1
	ChannelStatusActive
)

type ChannelConfig struct {
	ID             uuid.UUID `model:"role:pk,"`
	Name           string
	Node           *Node `model:"rel:belongs-to,join:NodeID=ID"`
	NodeID         int
	DataRate       telem.DataRate
	DataType       telem.DataType
	ConflictPolicy ChannelConflictPolicy
	Status         ChannelStatus
}

const MaxChunkSize = 2e7

type ChannelChunk struct {
	ID              uuid.UUID `model:"role:pk,"`
	Range           *Range
	RangeID         uuid.UUID
	ChannelConfig   *ChannelConfig `model:"rel:belongs-to,join:RangeID=ID"`
	ChannelConfigID uuid.UUID
	Size            int64
	StartTS         telem.TimeStamp
}

type ChannelChunkReplica struct {
	ID             uuid.UUID     `model:"role:pk,"`
	ChannelChunk   *ChannelChunk `model:"rel-belongs-to,join:ChannelChunkID=ID"`
	ChannelChunkID uuid.UUID
	RangeReplica   *RangeReplica `model:"rel:belongs-to,join:RangeReplicaID=ID"`
	RangeReplicaID uuid.UUID
	Telem          *telem.ChunkData `storage:"re:object," model:"role:telemChunkData,"`
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `model:"role:series"`
	ChannelConfigID uuid.UUID
	Value           float64         `storage:"role:cache"`
	Timestamp       telem.TimeStamp `storage:"role:cache"`
}
