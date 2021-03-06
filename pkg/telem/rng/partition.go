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
			Action:   pd.detectObserve,
			Interval: detectObserveInterval,
		},
	}
	defaultOpts := []tasks.ScheduleOpt{tasks.ScheduleWithName(scheduleWithName)}
	allOpts := append(defaultOpts, opts...)
	return tasks.NewScheduleSimple(tsk, allOpts...)

}

// |||| DETECT ||||

type partitionDetect struct {
	observe Observe
	qa      *QueryAssemble
}

func (pd *partitionDetect) detectObserve(ctx context.Context, opt tasks.ScheduleConfig) error {
	openRanges := pd.observe.RetrieveFilter(ObservedRange{Status: models.RangeStatusOpen})
	return pd.detect(ctx, openRanges, opt)
}

func (pd *partitionDetect) detectPersist(ctx context.Context, opt tasks.ScheduleConfig) error {
	var openRng []*models.Range
	if err := pd.qa.RetrieveOpenRangesQuery(&openRng).Exec(ctx); err != nil {
		return err
	}
	var obsRng []ObservedRange
	for _, openRange := range openRng {
		obsRng = append(obsRng, ObservedRange{
			PK:             openRange.ID,
			LeaseNodePK:    openRange.RangeLease.RangeReplica.NodeID,
			LeaseReplicaPK: openRange.RangeLease.RangeReplica.ID,
			Status:         openRange.Status,
		})
	}
	return pd.detect(ctx, obsRng, opt)
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
	pe := NewPartitionExecute(ctx, pd.qa, or.PK)
	oa, err := pe.OverAllocated()
	if !oa || err != nil {
		return []*models.Range{}, err
	}
	or.Status = models.RangeStatusPartition
	pd.observe.Add(or)
	newRng, err := pe.Exec()
	or.Status = models.RangeStatusClosed
	pd.observe.Add(or)
	return newRng, err
}

func (pd *partitionDetect) observeNewRngGroups(newRangeGroups [][]*models.Range) {
	for _, rangeGroup := range newRangeGroups {
		for _, rng := range rangeGroup {
			pd.observe.Add(ObservedRange{
				PK:             rng.ID,
				Status:         rng.Status,
				LeaseReplicaPK: rng.RangeLease.RangeReplica.ID,
				LeaseNodePK:    rng.RangeLease.RangeReplica.NodeID,
			})
		}
	}
}

// |||| EXECUTE ||||

func NewPartitionExecute(ctx context.Context, qa *QueryAssemble, rngPK uuid.UUID) *PartitionExecute {
	return &PartitionExecute{qa: qa, sourceRangePK: rngPK, catcher: errutil.NewCatchContext(ctx)}
}

// PartitionExecute checks if a models.Range is over-allocated (i.e. exceeds models.MaxRangeSize),
// and then splits it into smaller rngMap until its no longer allocated. It will then mark the original (source)
// models.Range as closed.
//
// Exec runs the partition, and returns any new models.Range objects created during the partition as well as
// any errors encountered.
type PartitionExecute struct {
	qa            *QueryAssemble
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
	newRng := p.createRange(sourceRng.RangeLease.RangeReplica.NodeID)
	p.reallocateChunks(excCC, newRng.ID)
	for nodePK, ccrPKC := range excCCR {
		newRRPK := p.newReplicaPK(newRng.RangeLease.RangeReplica, nodePK)
		p.reallocateChunkReplicas(ccrPKC, newRRPK)
	}
	p.newRanges = append(p.newRanges, newRng)
	p.updateRangeStatus(p.sourceRangePK, models.RangeStatusClosed)
	nextP := &PartitionExecute{qa: p.qa, sourceRangePK: newRng.ID, newRanges: p.newRanges, catcher: p.catcher}
	return nextP.Exec()
}

func (p *PartitionExecute) overAllocated() bool {
	rs := p.rangeSize()
	return rs > models.MaxRangeSize
}

func (p *PartitionExecute) rangeSize() (size int64) {
	p.catcher.Exec(p.qa.RetrieveRangeSizeQuery(p.sourceRangePK, &p._rngSize).Exec)
	return p._rngSize
}

func (p *PartitionExecute) updateRangeStatus(rngPK uuid.UUID, status models.RangeStatus) {
	p.catcher.Exec(p.qa.UpdateRangeStatusQuery(rngPK, status).Exec)
}

func (p *PartitionExecute) reallocateChunks(pks []uuid.UUID, rngPK uuid.UUID) {
	p.catcher.Exec(p.qa.ReallocateChunksQuery(pks, rngPK).Exec)
}

func (p *PartitionExecute) reallocateChunkReplicas(pks []uuid.UUID, rrPK uuid.UUID) {
	p.catcher.Exec(p.qa.ReallocateChunkReplicasQuery(pks, rrPK).Exec)
}

func (p *PartitionExecute) createRange(nodePK int) (newRng *models.Range) {
	p.catcher.Exec(func(ctx context.Context) (err error) {
		newRng, err = p.qa.CreateRange(ctx, nodePK)
		return err
	})
	return newRng
}

func (p *PartitionExecute) newReplicaPK(leaseRR *models.RangeReplica, nodeID int) uuid.UUID {
	newReplica := leaseRR
	if nodeID != leaseRR.NodeID {
		p.catcher.Exec(func(ctx context.Context) error {
			newRR, err := p.qa.CreateRangeReplica(ctx, leaseRR.RangeID, nodeID)
			newReplica = newRR
			return err
		})
	}
	return newReplica.ID
}

func (p *PartitionExecute) retrievePartitionInfo() (
	sourceRng *models.Range,
	sourceRR []*models.RangeReplica,
	cc []*models.ChannelChunk,
	ccr []*models.ChannelChunkReplica,
) {
	sourceRng = &models.Range{}
	p.catcher.Exec(p.qa.RetrieveRangeQuery(sourceRng, p.sourceRangePK).Exec)
	p.catcher.Exec(p.qa.RetrieveRangeChunksQuery(&cc, p.sourceRangePK).Exec)
	p.catcher.Exec(p.qa.RetrieveRangeReplicasQuery(&sourceRR, p.sourceRangePK).Exec)
	p.catcher.Exec(p.qa.RetrieveRangeChunkReplicasQuery(&ccr, p.sourceRangePK).Exec)
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
