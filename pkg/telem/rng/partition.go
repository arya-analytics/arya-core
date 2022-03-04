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
	scheduleWithName      = "Partition Scheduler"
	detectPersistInterval = 120 * time.Second
	detectObserveInterval = 30 * time.Second
)

func newSchedulePartition(pd *partitionDetect, opts ...tasks.ScheduleOpt) tasks.Schedule {
	tsk := []tasks.Task{
		{
			Name:     "Detect Persist",
			Action:   pd.detectPersist,
			Interval: detectPersistInterval,
		},
		{
			Name:     "Detect Observe",
			Action:   pd.detectObserver,
			Interval: detectObserveInterval,
		},
	}
	defaultOpts := []tasks.ScheduleOpt{tasks.ScheduleWithName(scheduleWithName)}
	allOpts := append(defaultOpts, opts...)
	return tasks.NewScheduleSimple(tsk, allOpts...)

}

// |||| DETECT ||||

type partitionDetect struct {
	Observe Observe
	Persist Persist
}

func (pd *partitionDetect) detectObserver(ctx context.Context, opt tasks.ScheduleConfig) error {
	openRanges := pd.Observe.RetrieveFilter(ObservedRange{Status: models.RangeStatusOpen})
	return pd.detect(ctx, openRanges, opt)
}

func (pd *partitionDetect) detectPersist(ctx context.Context, opt tasks.ScheduleConfig) error {
	openRanges, err := pd.Persist.RetrieveRangesByStatus(ctx)
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

func (pd *partitionDetect) detect(ctx context.Context, openRanges []ObservedRange, _ tasks.ScheduleConfig) error {
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

func (pd *partitionDetect) exec(ctx context.Context, or ObservedRange) ([]*models.Range, error) {
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

func (pd *partitionDetect) observeNewRngGroups(newRangeGroups [][]*models.Range) {
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
	return &PartitionExecute{pst: p, sourceRangePK: rngPK, catcher: errutil.NewCatchContext(ctx)}
}

// PartitionExecute checks if a models.Range is over-allocated (i.e. exceeds models.MaxRangeSize),
// and then splits it into smaller rngMap until its no longer allocated. It will then mark the original (source)
// models.Range as closed.
//
// Exec runs the partition, and returns any new models.Range objects created during the partition as well as
// any errors encountered.
type PartitionExecute struct {
	pst           Persist
	sourceRangePK uuid.UUID
	newRanges     []*models.Range
	catcher       *errutil.CatchContext
	_rngSize      int64
}

// OverAllocated returns true if the size of the source range exceeds models.MaxRangeSize.
func (p *PartitionExecute) OverAllocated() (bool, error) {
	return p.overAllocated(), p.catcher.Error()
}

// Exec executes the partition.
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
