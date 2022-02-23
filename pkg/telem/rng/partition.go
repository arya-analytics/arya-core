package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
)

type Partition struct {
	Persist   Persist
	RangeID   uuid.UUID
	NewRanges []*models.Range
}

func (p *Partition) Exec(ctx context.Context) ([]*models.Range, error) {
	size, err := p.Persist.RetrieveRangeSize(ctx, p.RangeID)
	if size < models.MaxRangeSize || err != nil {
		return p.NewRanges, nil
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
	for i := 0; i < len(chunks); i++ {
		if size < models.MaxRangeSize {
			break
		}
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

	reallocatedChunkReplicas := map[uuid.UUID][]uuid.UUID{}
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
			} else {
				reallocatedChunkReplicas[newRR.ID] = append(reallocatedChunkReplicas[newRR.ID], ccr.ID)
			}
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
	if err := p.Persist.CreateRangeReplica(ctx, &newRangeReplicas); err != nil {
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
	p.NewRanges = append(p.NewRanges, newRange)
	nextP := &Partition{Persist: p.Persist, RangeID: newRange.ID, NewRanges: p.NewRanges}
	return nextP.Exec(ctx)
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
