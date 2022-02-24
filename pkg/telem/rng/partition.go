package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/google/uuid"
	"sync"
	"time"
)

// |||| SCHEDULER ||||

const (
	detectPersistInterval = 120 * time.Second
	detectObserveInterval = 30 * time.Second
)

func NewSchedulePartition(pd *PartitionDetect, opts ...tasks.ScheduleOpt) tasks.Schedule {
	tsk := []tasks.Task{
		{
			Name:     "Detect Persist",
			Action:   pd.DetectPersist,
			Interval: detectPersistInterval,
		},
		{
			Name:     "Detect Observe",
			Action:   pd.DetectObserver,
			Interval: detectObserveInterval,
		},
	}
	defaultOpts := []tasks.ScheduleOpt{tasks.ScheduleWithName("Partition Scheduler")}
	allOpts := append(defaultOpts, opts...)
	return tasks.NewScheduleSimple(tsk, allOpts...)

}

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
	var or []ObservedRange
	for _, openRange := range openRanges {
		or = append(or, ObservedRange{
			PK:             openRange.ID,
			LeaseNodePK:    openRange.RangeLease.RangeReplica.NodeID,
			LeaseReplicaPK: openRange.RangeLease.RangeReplica.ID,
			Status:         openRange.Status,
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
	return &PartitionExecute{pst: p, sourceRangePK: rngPK, catcher: errutil.NewContextCatcher(ctx)}
}

type PartitionExecute struct {
	pst           Persist
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
		p._rngSize, err = p.pst.RetrieveRangeSize(ctx, p.sourceRangePK)
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
	for nodePK, ccrPKC := range excCCR {
		newRRPK := p.newReplicaPK(newRng.RangeLease.RangeReplica, nodePK)
		p.reallocateChunkReplicas(ccrPKC, newRRPK)
	}
	p.newRanges = append(p.newRanges, newRng)
	p.updateRangeStatus(p.sourceRangePK, models.RangeStatusClosed)
	nextP := &PartitionExecute{pst: p.pst, sourceRangePK: newRng.ID, newRanges: p.newRanges, catcher: p.catcher}
	return nextP.Exec()
}

func (p *PartitionExecute) updateRangeStatus(rngPK uuid.UUID, status models.RangeStatus) {
	p.catcher.Exec(func(ctx context.Context) error { return p.pst.UpdateRangeStatus(ctx, rngPK, status) })
}

func (p *PartitionExecute) reallocateChunks(ccPKs []uuid.UUID, rngPK uuid.UUID) {
	p.catcher.Exec(func(ctx context.Context) error { return p.pst.ReallocateChunks(ctx, ccPKs, rngPK) })
}

func (p *PartitionExecute) reallocateChunkReplicas(ccrPKC []uuid.UUID, RRPK uuid.UUID) {
	p.catcher.Exec(func(ctx context.Context) error { return p.pst.ReallocateChunkReplicas(ctx, ccrPKC, RRPK) })
}

func (p *PartitionExecute) newRange(nodePK int) (newRng *models.Range) {
	p.catcher.Exec(func(ctx context.Context) (err error) {
		newRng, err = p.pst.CreateRange(ctx, nodePK)
		return err
	})
	return newRng
}

func (p *PartitionExecute) newReplicaPK(leaseRR *models.RangeReplica, nodeID int) uuid.UUID {
	newReplicaID := leaseRR.ID
	if nodeID != leaseRR.NodeID {
		p.catcher.Exec(func(ctx context.Context) error {
			newRR, err := p.pst.CreateRangeReplica(ctx, leaseRR.RangeID, nodeID)
			newReplicaID = newRR.ID
			return err
		})
	}
	return newReplicaID
}

func (p *PartitionExecute) retrievePartitionInfo() (
	sourceRng *models.Range,
	sourceRR []*models.RangeReplica,
	cc []*models.ChannelChunk,
	ccr []*models.ChannelChunkReplica,
) {
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		sourceRng, cErr = p.pst.RetrieveRange(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		cc, cErr = p.pst.RetrieveRangeChunks(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		sourceRR, cErr = p.pst.RetrieveRangeReplicas(ctx, p.sourceRangePK)
		return cErr
	})
	p.catcher.Exec(func(ctx context.Context) (cErr error) {
		ccr, cErr = p.pst.RetrieveRangeChunkReplicas(ctx, p.sourceRangePK)
		return cErr
	})
	return sourceRng, sourceRR, cc, ccr
}

func excessChunkReplicas(rrC []*models.RangeReplica, ccrC []*models.ChannelChunkReplica, excCC []uuid.UUID) map[int][]uuid.UUID {
	excessCCR := map[int][]uuid.UUID{}
	for _, pk := range excCC {
		for _, ccr := range filterChunkReplicas(pk, ccrC) {
			rr, ok := findRangeReplica(ccr.RangeReplicaID, rrC)
			if !ok {
				panic("couldn't find the chunks range replica")
			}
			excessCCR[rr.NodeID] = append(excessCCR[rr.NodeID], ccr.ID)
		}
	}
	return excessCCR
}

func excessChunks(size int64, ccC []*models.ChannelChunk) (excCC []uuid.UUID) {
	for i := 0; i < len(ccC); i++ {
		if size < models.MaxRangeSize {
			break
		}
		c := ccC[i]
		excCC = append(excCC, c.ID)
		size -= c.Size
	}
	return excCC
}

func filterChunkReplicas(chunkPK uuid.UUID, CCR []*models.ChannelChunkReplica) (resCCR []*models.ChannelChunkReplica) {
	for _, ccr := range CCR {
		if ccr.ChannelChunkID == chunkPK {
			resCCR = append(resCCR, ccr)
		}
	}
	return resCCR
}

func findRangeReplica(RRPK uuid.UUID, RR []*models.RangeReplica) (*models.RangeReplica, bool) {
	for _, rr := range RR {
		if rr.ID == RRPK {
			return rr, true
		}
	}
	return nil, false
}
