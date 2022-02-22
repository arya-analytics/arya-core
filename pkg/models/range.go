package models

import "github.com/google/uuid"

const RangeSize int64 = 512e7

type RangeStatus int

const (
	RangeStatusOpen RangeStatus = iota + 1
	RangeStatusClosed
)

type Range struct {
	ID         uuid.UUID `model:"role:pk,"`
	Status     RangeStatus
	RangeLease *RangeLease
}

type RangeLease struct {
	ID             uuid.UUID `model:"role:pk,"`
	RangeID        uuid.UUID
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
