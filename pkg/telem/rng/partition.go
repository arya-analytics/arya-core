package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

// |||| DETECT ||||

type PartitionDetect struct {
	Observe Observe
	Persist Persist
}

func (pd *PartitionDetect) DetectObserver(ctx context.Context, opt tasks.ScheduleConfig) error {
	openRanges := pd.Observe.RetrieveFilter(ObservedRange{Status: models.RangeStatusOpen})
	return pd.detect(ctx, openRanges, opt)
}

func (pd *PartitionDetect) detect(ctx context.Context, openRanges []ObservedRange, opt tasks.ScheduleConfig) error {
	wg := sync.WaitGroup{}
	newRngGroups, errs := make([][]*models.Range, len(openRanges)), make([]error, len(openRanges))
	for i, or := range openRanges {
		wg.Add(1)
		go func(i int, or ObservedRange) {
			newRngGroups[i], errs[i] = pd.exec(ctx, or)
			wg.Done()
		}(i, or)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	pd.observeNewRngGroups(newRngGroups)
	return nil

}

func (pd *PartitionDetect) DetectPersist(ctx context.Context, opt tasks.ScheduleConfig) error {
	openRanges, err := pd.Persist.RetrieveOpenRanges(ctx)
	if err != nil {
		return err
	}
	log.Info(len(openRanges))
	var or []ObservedRange
	for _, openRange := range openRanges {
		or = append(or, ObservedRange{
			PK:             openRange.ID,
			LeaseNodePK:    openRange.RangeLease.RangeReplica.NodeID,
			LeaseReplicaPK: openRange.RangeLease.RangeReplica.ID,
		})
	}
	return pd.detect(ctx, or, opt)
}

func (pd *PartitionDetect) exec(ctx context.Context, or ObservedRange) ([]*models.Range, error) {
	pe := NewPartitionExecute(ctx, pd.Persist, or.PK)
	oa, err := pe.OverAllocated()
	if !oa || err != nil {
		return []*models.Range{}, err
	}
	or.Status = models.RangeStatusPartition
	pd.Observe.Add(or)
	newRng, err := pe.Exec()
	or.Status = models.RangeStatusClosed
	pd.Observe.Add(or)
	return newRng, err
}

func (pd *PartitionDetect) observeNewRngGroups(newRangeGroups [][]*models.Range) {
	for _, rangeGroup := range newRangeGroups {
		for _, rng := range rangeGroup {
			pd.Observe.Add(ObservedRange{
				PK:             rng.ID,
				Status:         rng.Status,
				LeaseReplicaPK: rng.RangeLease.RangeReplica.ID,
				LeaseNodePK:    rng.RangeLease.RangeReplica.NodeID,
			})
		}
	}
}

// |||| EXECUTE ||||

func NewPartitionExecute(ctx context.Context, p Persist, rngPK uuid.UUID) *PartitionExecute {
	return &PartitionExecute{Per: p, sourceRangePK: rngPK, catcher: errutil.NewContextCatcher(ctx)}
}

type PartitionExecute struct {
	Per           Persist
	sourceRangePK uuid.UUID
	newRanges     []*models.Range
	catcher       *errutil.ContextCatcher
	_rngSize      int64
}

func (p *PartitionExecute) OverAllocated() (bool, error) {
	return p.overAllocated(), p.catcher.Error()
}

func (p *PartitionExecute) overAllocated() bool {
	return p.rangeSize() > models.MaxRangeSize
}

func (p *PartitionExecute) rangeSize() (size int64) {
	p.catcher.Exec(func(ctx context.Context) (err error) {
		p._rngSize, err = p.Per.RetrieveRangeSize(ctx, p.sourceRangePK)
		return err
	})
	return p._rngSize
}

func (p *PartitionExecute) Exec() ([]*models.Range, error) {
	if !p.overAllocated() {
		return p.newRanges, nil
	}
	p.updateRangeStatus(p.sourceRangePK, models.RangeStatusPartition)
	sourceRng, sourceRR, cc, ccr := p.retrievePartitionInfo()
	excCC := excessChunks(p.rangeSize(), cc)
	excCCR := excessChunkReplicas(sourceRR, ccr, excCC)
	newRng := p.newRange(sourceRng.RangeLease.RangeReplica.NodeID)
	p.reallocateChunks(excCC, newRng.ID)
	for nodePK, ccrPKs := range excCCR {
		newRRPK := p.newReplicaPK(newRng.RangeLease.RangeReplica, nodePK)
		p.reallocateChunkReplicas(ccrPKs, newRRPK)
	}
	p.newRanges = append(p.newRanges, newRng)
	p.updateRangeStatus(p.sourceRangePK, models.RangeStatusClosed)
	nextP := &PartitionExecute{Per: p.Per, sourceRangePK: newRng.ID, newRanges: p.newRanges, catcher: p.catcher}
	return nextP.Exec()
}

func (p *PartitionExecute) updateRangeStatus(rngPK uuid.UUID, status models.RangeStatus) {
	p.catcher.Exec(func(ctx context.Context) error { return p.Per.UpdateRangeStatus(ctx, rngPK, status) })
}

func (p *PartitionExecute) reallocateChunks(ccPKs []uuid.UUID, rngPK uuid.UUID) {
	p.catcher.Exec(func(ctx context.Context) error { return p.Per.ReallocateChunks(ctx, ccPKs, rngPK) })
}

func (p *PartitionExecute) reallocateChunkReplicas(ccrPKs []uuid.UUID, RRPK uuid.UUID) {
	p.catcher.Exec(func(ctx context.Context) error { return p.Per.ReallocateChunkReplicas(ctx, ccrPKs, RRPK) })
}

func (p *PartitionExecute) newRange(nodePK int) (newRng *models.Range) {
	p.catcher.Exec(func(ctx context.Context) (err error) {
		newRng, err = p.Per.NewRange(ctx, nodePK)
		return err
	})
	return newRng
}

func (p *PartitionExecute) newReplicaPK(leaseRR *models.RangeReplica, nodeID int) uuid.UUID {
	newReplicaID := leaseRR.ID
	if nodeID != leaseRR.NodeID {
		p.catcher.Exec(func(ctx context.Context) error {
			newRR, err := p.Per.NewRangeReplica(ctx, leaseRR.RangeID, nodeID)
			newReplicaID = newRR.ID
			return err
		})
	}
	return newReplicaID
}

func (p *PartitionExecute) retrievePartitionInfo() (sourceRng *models.Range, sourceRR []*models.RangeReplica, cc []*models.ChannelChunk, ccr []*models.ChannelChunkReplica) {
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		sourceRng, cErr = p.Per.RetrieveRange(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		cc, cErr = p.Per.RetrieveRangeChunks(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		sourceRR, cErr = p.Per.RetrieveRangeReplicas(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		ccr, cErr = p.Per.RetrieveRangeChunkReplicas(ctx, p.sourceRangePK)
		return cErr
	})
	return sourceRng, sourceRR, cc, ccr
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
