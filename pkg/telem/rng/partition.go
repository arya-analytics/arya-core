package rng

//type PartitionScheduler struct {
//	s Service
//	*tasks.SchedulerSimple
//}
//
//func NewPartitionScheduler(observerCheckInterval, clusterCheckInterval time.Duration) *PartitionScheduler {
//	return &PartitionScheduler{}
//
//}
//
//func (ps *PartitionScheduler) scanObservedRanges(ctx context.Context, cfg tasks.SchedulerConfig) error {
//	openRanges := ps.s.obs.RetrieveFilter(ObservedRange{Status: models.RangeStatusOpen})
//	var wg = sync.WaitGroup{}
//	var err error
//	for _, or := range openRanges {
//		go func(or ObservedRange) {
//			or.Status = models.RangeStatusPartition
//			ps.s.obs.Add(or)
//			newRanges, pErr := NewPartition(ps.p, or.ID).Exec(ctx)
//			if pErr != nil {
//				err = pErr
//			}
//			or.Status = models.RangeStatusClosed
//			ps.s.obs.Add(or)
//			for _, rng := range newRanges {
//				ps.s.obs.Add(ObservedRange{
//					ID:             rng.ID,
//					LeaseNodeID:    rng.RangeLease.RangeReplica.NodeID,
//					LeaseReplicaID: rng.RangeLease.RangeReplica.ID,
//					Status:         rng.Status,
//				})
//			}
//			wg.Done()
//		}(or)
//	}
//	wg.Wait()
//	return err
//}
//
//const unCalculatedRangeSize = -1
//
//func NewPartition(p Persist, rangeID uuid.UUID) *Partition {
//	return &Partition{sourceRangeID: rangeID, _rangeSize: unCalculatedRangeSize, catcher: &errutil.Catcher{}}
//}
//
//type Partition struct {
//	persist       Persist
//	ctx           context.Context
//	sourceRangeID uuid.UUID
//	catcher       *errutil.Catcher
//	_sourceRange  *models.Range
//	_rangeSize    int64
//}
//
//func (p *Partition) Exec(ctx context.Context) ([]*models.Range, error) {
//	p.ctx = ctx
//	if p.rangeSize() < models.MaxRangeSize {
//		return []*models.Range{}, nil
//	}
//	newRanges := p.createRanges(p.rangeSize(), p.sourceRange().RangeLease.RangeReplica.NodeID)
//	replicas := p.retrieveChunkReplicas()
//	reAllocatedReplicas []*models.RangeReplica{}
//	for _, rng := range newRanges {
//
//	}
//	return newRanges, p.catcher.Error()
//}
//
//func (p *Partition) allocateUntilFull(replicas *models.RangeReplica, reAllocated) {
//
//}
//
//func (p *Partition) retrieveChunkReplicas() (replicas []*models.ChannelChunkReplica) {
//	p.catcher.Exec(func() error {
//		return p.persist.RetrieveChunkReplicas(p.ctx, replicas, p.sourceRange().RangeLease.RangeReplica.ID)
//	})
//	return replicas
//}
//
//func (p *Partition) createRanges(rangeSize int64, nodeID int) (newRanges []*models.Range) {
//	count := newRangeCount(rangeSize)
//	p.catcher.Exec(func() (err error) {
//		newRanges, err = p.persist.CreateRanges(p.ctx, nodeID, count)
//		return err
//	})
//	return newRanges
//}
//
//func (p *Partition) rangeSize() int64 {
//	if p._rangeSize == unCalculatedRangeSize {
//		p.catcher.Exec(func() (err error) {
//			p._rangeSize, err = p.persist.RetrieveRangeSize(p.ctx, p.sourceRange().ID)
//			return err
//		})
//	}
//	return p._rangeSize
//}
//
//func (p *Partition) sourceRange() *models.Range {
//	if p._sourceRange == nil {
//		p.catcher.Exec(func() error {
//			return p.persist.RetrieveRange(p.ctx, p._sourceRange, p.sourceRangeID)
//		})
//	}
//	return p._sourceRange
//}
//
//func newRangeCount(rangeSize int64) int {
//	return int(math.Ceil(float64(rangeSize)/float64(models.MaxRangeSize) - 1))
//}

/* WHAT DOES THE PARTITION LOGIC LOOK LIKE

1. At some arbitrary interval
2. Check for overallocated ranges
3. If the range is overallocated
	1. Mark the range as RangeStatusPartition
	2. Partition the range.
	3. Mark	the new range as open and add it to observer
	4. Mark the full range as closed and remote it from the  observer
*/
