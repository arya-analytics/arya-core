package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

// |||| BASE ALLOCATOR ||||

type Allocate struct {
	p       Persist
	obs     Observe
	nodeID  int
	rangeID uuid.UUID
}

func (a *Allocate) Chunk(nodeID int, chunk *models.ChannelChunk) *allocateChunk {
	a.nodeID = nodeID
	return &allocateChunk{Allocate: a, chunk: chunk}
}

func (a *Allocate) ChunkReplica(replica *models.ChannelChunkReplica) *allocateChunkReplica {
	if a.nodeID == 0 || model.NewPK(a.rangeID).IsZero() {
		panic("can't allocate a chunk replica before a chunk")
	}
	return &allocateChunkReplica{Allocate: a, replica: replica}
}

func (a *Allocate) retrieveObservedOrNew(ctx context.Context, q ObservedRange) (ObservedRange, error) {
	or, ok := a.obs.Retrieve(q)
	if !ok {
		var err error
		newRng, err := a.p.NewRange(ctx, a.nodeID)
		if err != nil {
			return ObservedRange{}, err
		}
		or = ObservedRange{
			ID:             newRng.ID,
			Status:         newRng.Status,
			LeaseReplicaID: newRng.RangeLease.RangeReplica.ID,
			LeaseNodeID:    newRng.RangeLease.RangeReplica.NodeID,
		}
		a.obs.Add(or)
	}
	return or, nil
}

// |||| CHUNK ALLOCATOR ||||

type allocateChunk struct {
	*Allocate
	chunk *models.ChannelChunk
}

func (ac *allocateChunk) Exec(ctx context.Context) error {
	or, err := ac.retrieveObservedOrNew(ctx, ObservedRange{LeaseNodeID: ac.nodeID, Status: models.RangeStatusOpen})
	ac.rangeID = or.ID
	ac.chunk.RangeID = or.ID
	return err
}

// |||| CHUNK REPLICA ALLOCATOR ||||

type allocateChunkReplica struct {
	*Allocate
	replica *models.ChannelChunkReplica
}

func (ac *allocateChunkReplica) Exec(ctx context.Context) error {
	or, err := ac.retrieveObservedOrNew(ctx, ObservedRange{ID: ac.rangeID, Status: models.RangeStatusOpen})
	ac.replica.RangeReplicaID = or.LeaseReplicaID
	return err
}
