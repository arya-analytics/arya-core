package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
)

// |||| BASE ALLOCATOR ||||

// Allocate allocates ChannelChunks to Ranges and their child ChannelChunkReplicas to RangeReplicas.
// Allocate should not be instantiated directly, and should instead be instantiated through
// rng.Service by calling NewAllocate.
//
// A new Allocate should be created for every unique ChannelChunk.
type Allocate struct {
	obs         Observe
	exec        query.Execute
	leaseNodePK int
	rangePK     uuid.UUID
}

// Chunk creates a new allocateChunk struct which can be used to allocate a ChannelChunk to
// a Range.
func (a *Allocate) Chunk(leaseNodePK int, chunk *models.ChannelChunk) *allocateChunk {
	a.leaseNodePK = leaseNodePK
	return &allocateChunk{Allocate: a, chunk: chunk}
}

// ChunkReplica creates a new allocateChunkReplica which can be used to allocate a ChannelChunkReplica
// to a Range.
//
// NOTE: Allocate.Chunk must be called and executed before calling this method.
func (a *Allocate) ChunkReplica(replica *models.ChannelChunkReplica) *allocateChunkReplica {
	return &allocateChunkReplica{Allocate: a, replica: replica}
}

func (a *Allocate) retrieveObservedOrCreate(ctx context.Context, q ObservedRange) (ObservedRange, error) {
	or, ok := a.obs.Retrieve(q)
	if !ok {
		newRng, err := createRange(ctx, a.exec, a.leaseNodePK)
		if err != nil {
			return ObservedRange{}, err
		}
		or = ObservedRange{
			PK:             newRng.ID,
			Status:         newRng.Status,
			LeaseReplicaPK: newRng.RangeLease.RangeReplica.ID,
			LeaseNodePK:    newRng.RangeLease.RangeReplica.NodeID,
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

// Exec runs the allocator, and assigns the appropriate RangeID to the ChannelChunk.
func (ac *allocateChunk) Exec(ctx context.Context) error {

	or, err := ac.retrieveObservedOrCreate(ctx, ObservedRange{LeaseNodePK: ac.leaseNodePK, Status: models.RangeStatusOpen})
	ac.rangePK = or.PK
	ac.chunk.RangeID = or.PK
	return err
}

// |||| CHUNK REPLICA ALLOCATOR ||||

type allocateChunkReplica struct {
	*Allocate
	replica *models.ChannelChunkReplica
}

// Exec runs the allocator, and assigns the appropriate RangeReplicaID to the ChannelChunkReplica.
func (ac *allocateChunkReplica) Exec(ctx context.Context) error {
	if ac.leaseNodePK == 0 || model.NewPK(ac.rangePK).IsZero() {
		panic("can't allocate a chunk replica before a chunk")
	}
	or, err := ac.retrieveObservedOrCreate(ctx, ObservedRange{PK: ac.rangePK, Status: models.RangeStatusOpen})
	ac.replica.RangeReplicaID = or.LeaseReplicaPK
	return err
}
