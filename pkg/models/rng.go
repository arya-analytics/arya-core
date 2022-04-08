package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

const MaxRangeSize int64 = 64e7

type RangeStatus int

const (
	RangeStatusClosed RangeStatus = iota + 1
	RangeStatusOpen
	RangeStatusPartition
)

type Range struct {
	model.Base `storage:"engines:md,"`
	ID         uuid.UUID   `bun:"type:UUID,pk" model:"role:pk,"`
	Status     RangeStatus `bun:"type:SMALLINT,default:2"`
	Open       bool        `bun:"default:TRUE,"`
	RangeLease *RangeLease `model:"rel:has-one,join:ID=RangeID" bun:"rel:has-one,join:id=range_id"`
}

type RangeLease struct {
	model.Base     `storage:"engines:md,"`
	ID             uuid.UUID     `bun:"type:UUID,pk" model:"role:pk,"`
	RangeID        uuid.UUID     `bun:"type:UUID"`
	RangeReplica   *RangeReplica `bun:"rel:belongs-to,join:range_replica_id=id" model:"rel:belongs-to,join:RangeReplicaID=ID"`
	RangeReplicaID uuid.UUID     `model:"rel:belongs-to,join:RangeReplicaID=ID" bun:"type:UUID,"`
}

type RangeReplica struct {
	model.Base `storage:"engines:md,"`
	ID         uuid.UUID `bun:"type:UUID,pk" model:"role:pk,"`
	Range      *Range    `model:"rel:belongs-to,join:RangeID=ID" bun:"rel:belongs-to,join:range_id=id"`
	RangeID    uuid.UUID `bun:"type:UUID,"`
	Node       *Node     `model:"rel:belongs-to,join:NodeID=ID" bun:"rel:belongs-to,join:node_id=id"`
	NodeID     int
}
