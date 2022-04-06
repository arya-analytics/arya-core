package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
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
	model.Base     `model:"role:tsSeries" storage:"engines:md+cache," `
	ID             uuid.UUID `model:"role:pk," bun:"type:UUID,pk"`
	Name           string
	Node           *Node `model:"rel:belongs-to,join:NodeID=ID" bun:"rel:belongs-to,join:node_id=id,"`
	NodeID         int
	DataRate       telem.DataRate        `bun:"type:decimal"`
	DataType       telem.DataType        `bun:"default:0"`
	ConflictPolicy ChannelConflictPolicy `bun:"default:1"`
	Status         ChannelStatus         `bun:"default:1"`
	Retention      telem.TimeSpan        `bun:"type:int64"`
}

const MaxChunkSize = 2e7

type ChannelChunk struct {
	model.Base      `storage:"engines:md,"`
	ID              uuid.UUID      `model:"role:pk," bun:"type:UUID,pk"`
	Range           *Range         `bun:"rel:belongs-to,join:range_id=id,"`
	RangeID         uuid.UUID      `bun:"type:UUID,"`
	ChannelConfig   *ChannelConfig `model:"rel:belongs-to,join:RangeID=ID" bun:"rel:belongs-to,join:channel_config_id=id,"`
	ChannelConfigID uuid.UUID      `bun:"type:UUID,"`
	Size            int64
	StartTS         telem.TimeStamp `bun:"type:int64"`
}

type ChannelChunkReplica struct {
	model.Base     `storage:"engines:md+obj,"`
	ID             uuid.UUID        `bun:"type:UUID,pk" model:"role:pk,"`
	ChannelChunk   *ChannelChunk    `model:"rel:belongs-to,join:ChannelChunkID=ID" bun:"rel:belongs-to,join:channel_chunk_id=id,"`
	ChannelChunkID uuid.UUID        `bun:"type:UUID,"`
	RangeReplica   *RangeReplica    `model:"rel:belongs-to,join:RangeReplicaID=ID" bun:"rel:belongs-to,join:range_replica_id=id,"`
	RangeReplicaID uuid.UUID        `bun:"type:UUID,"`
	Telem          *telem.ChunkData `storage:"re:object," model:"role:telemChunkData,"`
}

type ChannelSample struct {
	model.Base      `model:"role:tsSample" storage:"engines:cache,"`
	ChannelConfig   *ChannelConfig  `model:"role:series"`
	ChannelConfigID uuid.UUID       `model:"role:pk,"`
	Value           float64         `model:"role:tsValue," storage:"role:cache"`
	Timestamp       telem.TimeStamp `model:"role:tsStamp," storage:"role:cache"`
}
