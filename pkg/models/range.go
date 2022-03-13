package models

import "github.com/google/uuid"

const MaxRangeSize int64 = 512e7

type RangeStatus int

const (
	RangeStatusClosed RangeStatus = iota + 1
	RangeStatusOpen
	RangeStatusPartition
)

type Range struct {
	ID         uuid.UUID `model:"role:pk,"`
	Status     RangeStatus
	RangeLease *RangeLease `model:"rel:has-one,join:ID=RangeID"`
}

type RangeLease struct {
	ID             uuid.UUID `model:"role:pk,"`
	RangeID        uuid.UUID
	RangeReplica   *RangeReplica `model:"rel:belongs-to,join:RangeReplicaID=ID"`
	RangeReplicaID uuid.UUID
}

type RangeReplica struct {
	ID      uuid.UUID `model:"role:pk"`
	Range   *Range    `model:"rel:belongs-to,join:RangeID=ID"`
	RangeID uuid.UUID
	Node    *Node `model:"rel:belongs-to,join:NodeID=ID"`
	NodeID  int
}
