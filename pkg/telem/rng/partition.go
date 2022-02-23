package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
)

type Partition struct {
	Persist Persist
	RangeID uuid.UUID
}

/* Order of opts

1. Retrieve the size of the leaseholder range replica

// If we need to partition

2. Retrieve all chunks belonging to the leaseholder range

3. For all chunks
	1. Reassign the chunk to the new range

3. Retrieve all chunk replicas belonging to the reassigned chunk
	2. Reassing teh chunk replicas to the new range replicas

4. save everything and return values
*/

func (p *Partition) Exec(ctx context.Context) (*models.Range, error) {
	size, err := p.Persist.RetrieveRangeSize(ctx, p.RangeID)
	if size < models.MaxRangeSize || err != nil {
		return nil, nil
	}

	sourceRange, err := p.Persist.RetrieveRange(ctx, p.RangeID)
	if err != nil {
		return nil, err
	}

	chunks, err := p.Persist.RetrieveRangeChunks(ctx, p.RangeID)
	if err != nil {
		return nil, err
	}

	var reallocatedChunkIDS []uuid.UUID
	for i := 0; size > models.MaxRangeSize; i++ {
		c := chunks[i]
		reallocatedChunkIDS = append(reallocatedChunkIDS, c.ID)
		size -= c.Size
	}

	newRange := &models.Range{ID: uuid.New()}
	var newRangeReplicas []*models.RangeReplica

	rangeReplicas, err := p.Persist.RetrieveRangeReplicas(ctx, p.RangeID)
	if err != nil {
		return nil, err
	}

	chunkReplicas, err := p.Persist.RetrieveRangeChunkReplicas(ctx, p.RangeID)
	if err != nil {
		return nil, err
	}

	var reallocatedChunkReplicas map[uuid.UUID][]uuid.UUID
	for _, ID := range reallocatedChunkIDS {
		ccrs := findChunkReplicas(ID, chunkReplicas)
		for _, ccr := range ccrs {
			rr, ok := findRangeReplica(ccr.RangeReplicaID, rangeReplicas)
			if !ok {
				panic("couldn't find the chunks range replica")
			}
			newRR, ok := findRangeReplicaByNodeID(rr.NodeID, newRangeReplicas)
			if !ok {
				newRR = &models.RangeReplica{ID: uuid.New(), RangeID: newRange.ID, NodeID: rr.NodeID}
				newRangeReplicas = append(newRangeReplicas, newRR)
				reallocatedChunkReplicas[newRR.ID] = []uuid.UUID{ccr.ID}
			}
			reallocatedChunkReplicas[newRR.ID] = append(reallocatedChunkReplicas[newRR.ID], ccr.ID)
		}
	}

	newLeaseReplica, ok := findRangeReplicaByNodeID(sourceRange.RangeLease.RangeReplica.NodeID, newRangeReplicas)
	if !ok {
		panic("couldn't find new lease replica")
	}
	newLease := &models.RangeLease{
		ID:             uuid.New(),
		RangeID:        newRange.ID,
		RangeReplicaID: newLeaseReplica.ID,
	}

	if err := p.Persist.CreateRange(ctx, newRange); err != nil {
		return nil, err
	}
	if err := p.Persist.CreateRangeReplica(ctx, newRangeReplicas); err != nil {
		return nil, err
	}
	if err := p.Persist.CreateRangeLease(ctx, newLease); err != nil {
		return nil, err
	}
	if err := p.Persist.ReallocateChunks(ctx, reallocatedChunkIDS, newRange.ID); err != nil {
		return nil, err
	}
	for replicaID, chunkReplicaIDS := range reallocatedChunkReplicas {
		if err := p.Persist.ReallocateChunkReplicas(ctx, chunkReplicaIDS, replicaID); err != nil {
			return nil, err
		}
	}

	newLease.RangeReplica = newLeaseReplica
	newRange.RangeLease = newLease
	return newRange, nil
}

func (p *Partition) rangeExceedsMaxSize(ctx context.Context) (bool, error) {
	size, err := p.Persist.RetrieveRangeSize(ctx, p.RangeID)
	return size > models.MaxRangeSize, err
}

func findChunkReplicas(chunkID uuid.UUID, chunkReplicas []*models.ChannelChunkReplica) (results []*models.ChannelChunkReplica) {
	for _, ccr := range chunkReplicas {
		if ccr.ChannelChunkID == chunkID {
			results = append(results, ccr)
		}
	}
	return results
}

func findRangeReplica(rangeReplicaID uuid.UUID, rangeReplicas []*models.RangeReplica) (*models.RangeReplica, bool) {
	for _, rr := range rangeReplicas {
		if rr.ID == rangeReplicaID {
			return rr, true
		}
	}
	return nil, false
}

func findRangeReplicaByNodeID(nodeID int, rangeReplicas []*models.RangeReplica) (*models.RangeReplica, bool) {
	for _, rr := range rangeReplicas {
		if rr.NodeID == nodeID {
			return rr, true
		}
	}
	return nil, false

}
