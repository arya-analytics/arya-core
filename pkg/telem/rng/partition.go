package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
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

	sourceRng, sourceRR, cc, ccr, err := p.retrievePartitionInfo(ctx)
	if err != nil {
		return nil, err
	}

	excCC := excessChunks(size, cc)
	excCCR := excessChunkReplicas(sourceRR, ccr, excCC)

	newRange, err := p.Persist.NewRange(ctx, sourceRng.RangeLease.RangeReplica.NodeID)
	if err != nil {
		return nil, err
	}

	if cErr := p.Persist.ReallocateChunks(ctx, excCC, newRange.ID); cErr != nil {
		return nil, cErr
	}

	for nodeID, chunkReplicaIDS := range excCCR {
		newReplicaID := newRange.RangeLease.RangeReplicaID
		if nodeID != newRange.RangeLease.RangeReplica.NodeID {
			newRR, nRRErr := p.Persist.NewRangeReplica(ctx, newRange.ID, nodeID)
			if nRRErr != nil {
				return nil, nRRErr
			}
			newReplicaID = newRR.ID
		}
		if ccrErr := p.Persist.ReallocateChunkReplicas(ctx, chunkReplicaIDS, newReplicaID); ccrErr != nil {
			return nil, ccrErr
		}
	}

	p.NewRanges = append(p.NewRanges, newRange)
	nextP := &Partition{Persist: p.Persist, RangeID: newRange.ID, NewRanges: p.NewRanges}
	return nextP.Exec(ctx)
}

func (p *Partition) retrievePartitionInfo(ctx context.Context) (sourceRng *models.Range,
	sourceRR []*models.RangeReplica,
	cc []*models.ChannelChunk,
	ccr []*models.ChannelChunkReplica, err error) {
	c := &errutil.Catcher{}
	c.Exec(func() (cErr error) {
		sourceRng, cErr = p.Persist.RetrieveRange(ctx, p.RangeID)
		return cErr
	})
	c.Exec(func() (cErr error) {
		cc, cErr = p.Persist.RetrieveRangeChunks(ctx, p.RangeID)
		return cErr
	})
	c.Exec(func() (cErr error) {
		sourceRR, cErr = p.Persist.RetrieveRangeReplicas(ctx, p.RangeID)
		return cErr
	})
	c.Exec(func() (cErr error) {
		ccr, cErr = p.Persist.RetrieveRangeChunkReplicas(ctx, p.RangeID)
		return cErr
	})
	return sourceRng, sourceRR, cc, ccr, c.Error()
}

func excessChunkReplicas(
	rrC []*models.RangeReplica,
	ccrC []*models.ChannelChunkReplica,
	excessCC []uuid.UUID) map[int][]uuid.UUID {
	excessCCR := map[int][]uuid.UUID{}
	for _, ID := range excessCC {
		for _, ccr := range findChunkReplicas(ID, ccrC) {
			rr, ok := findRangeReplica(ccr.RangeReplicaID, rrC)
			if !ok {
				panic("couldn't find the chunks range replica")
			}
			excessCCR[rr.NodeID] = append(excessCCR[rr.NodeID], ccr.ID)
		}
	}
	return excessCCR
}

func excessChunks(size int64, chunks []*models.ChannelChunk) (reallocatedChunkIDs []uuid.UUID) {
	for i := 0; i < len(chunks); i++ {
		if size < models.MaxRangeSize {
			break
		}
		c := chunks[i]
		reallocatedChunkIDs = append(reallocatedChunkIDs, c.ID)
		size -= c.Size
	}
	return reallocatedChunkIDs
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
