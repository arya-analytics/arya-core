package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

type ChannelConfig struct {
	ID       uuid.UUID `model:"role:pk,"`
	Name     string
	Node     *Node
	NodeID   int
	DataRate int
}

const MaxChunkSize = 2e7

type ChannelChunk struct {
	ID              uuid.UUID `model:"role:pk,"`
	Range           *Range
	RangeID         uuid.UUID
	ChannelConfig   *ChannelConfig
	ChannelConfigID uuid.UUID
	Size            int64
	StartTS         int64
}

type ChannelChunkReplica struct {
	ID             uuid.UUID `model:"role:pk,"`
	ChannelChunk   *ChannelChunk
	ChannelChunkID uuid.UUID
	RangeReplica   *RangeReplica
	RangeReplicaID uuid.UUID
	Telem          *telem.Bulk `storage:"re:object," model:"role:bulkTelem,"`
}

type ChannelSample struct {
	ChannelConfig   *ChannelConfig `model:"role:series"`
	ChannelConfigID uuid.UUID
	Value           float64 `storage:"role:cache"`
	Timestamp       int64   `storage:"role:cache"`
}
