package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

type ChannelConflictPolicy int

//go:generate stringer -type=ChannelConflictPolicy
const (
	ChannelConflictPolicyError ChannelConflictPolicy = iota + 1
	ChannelConflictPolicyDiscard
	ChannelConflictPolicyOverwrite
)

type ChannelState int

//go:generate stringer -type=ChannelState
const (
	ChannelStateInactive ChannelState = iota + 1
	ChannelStateActive
)

type ChannelConfig struct {
	ID             uuid.UUID `model:"role:pk,"`
	Name           string
	Node           *Node
	NodeID         int
	DataRate       telem.DataRate
	DataType       telem.DataType
	ConflictPolicy ChannelConflictPolicy `bun:"default:1"`
	State          ChannelState          `bun:"default:1"`
}

const MaxChunkSize = 2e7

type ChannelChunk struct {
	ID              uuid.UUID `model:"role:pk,"`
	Range           *Range
	RangeID         uuid.UUID
	ChannelConfig   *ChannelConfig
	ChannelConfigID uuid.UUID
	Size            int64
	StartTS         telem.TimeStamp
}

type ChannelChunkReplica struct {
	ID             uuid.UUID `model:"role:pk,"`
	ChannelChunk   *ChannelChunk
	ChannelChunkID uuid.UUID
	RangeReplica   *RangeReplica
	RangeReplicaID uuid.UUID
	Telem          *telem.ChunkData `storage:"re:object," model:"role:telemChunkData,"`
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `model:"role:series"`
	ChannelConfigID uuid.UUID
	Value           float64         `storage:"role:cache"`
	Timestamp       telem.TimeStamp `storage:"role:cache"`
}
