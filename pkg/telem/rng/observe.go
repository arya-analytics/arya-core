package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	"reflect"
	"sync"
)

// Observe watches a set of models.Range in an in-memory
// store in order to provide fast, goroutine secure retrieve
// and update operations.
//
// Observe should not persist any changes to permanent storage.
type Observe interface {
	// Add adds a new ObservedRange to Observe. Updates an ObservedRange if the ObservedRange.PK
	// already exists in Observe.
	//
	// NOTE: ObservedRange must be completely defined (i.e. all fields nonzero), or else the method will panic.
	Add(or ObservedRange)
	// Retrieve retrieves an ObservedRange matching a (partially) defined query q. Returns false if the
	// ObservedRange can't be found.
	Retrieve(q ObservedRange) (ObservedRange, bool)
	// RetrieveAll retrieves all ObservedRange currently in Observe.
	RetrieveAll() []ObservedRange
	// RetrieveFilter retrieves all ObservedRange matching a (partially) defined query q.
	RetrieveFilter(q ObservedRange) []ObservedRange
}

// ObservedRange holds critical info about a tracked range.
type ObservedRange struct {
	PK             uuid.UUID
	Status         models.RangeStatus
	LeaseNodePK    int
	LeaseReplicaPK uuid.UUID
}

// ObserveMem implements Observe with a mutex controlled in-memory map.
type ObserveMem struct {
	mu     sync.Mutex
	rngMap map[uuid.UUID]ObservedRange
}

// NewObserveMem creates a new ObserveMem with the preloaded rngMap.
func NewObserveMem(ranges []ObservedRange) *ObserveMem {
	rngMap := map[uuid.UUID]ObservedRange{}
	for _, rng := range ranges {
		validateObservedRange(rng)
		rngMap[rng.PK] = rng
	}
	return &ObserveMem{rngMap: rngMap}
}

func (o *ObserveMem) Add(or ObservedRange) {
	o.mu.Lock()
	defer o.mu.Unlock()
	validateObservedRange(or)
	o.rngMap[or.PK] = or
}

func (o *ObserveMem) Retrieve(q ObservedRange) (ObservedRange, bool) {
	matches := o.RetrieveFilter(q)
	if len(matches) > 0 {
		return matches[0], true
	}
	return ObservedRange{}, false
}

func (o *ObserveMem) RetrieveAll() (ranges []ObservedRange) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// Copying here to protect against modification
	for _, v := range o.rngMap {
		ranges = append(ranges, v)
	}
	return ranges
}

func (o *ObserveMem) RetrieveFilter(q ObservedRange) (matches []ObservedRange) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, or := range o.rngMap {
		if matchObservedRangeQuery(q, or) {
			matches = append(matches, or)
		}
	}
	return matches
}

func matchObservedRangeQuery(q ObservedRange, or ObservedRange) bool {
	if q.Status != 0 && q.Status != or.Status {
		return false
	}
	if !model.NewPK(q.PK).IsZero() && q.PK != or.PK {
		return false
	}
	if !model.NewPK(q.LeaseReplicaPK).IsZero() && q.LeaseReplicaPK != or.LeaseReplicaPK {
		return false
	}
	if !model.NewPK(q.LeaseNodePK).IsZero() && q.LeaseNodePK != or.LeaseNodePK {
		return false
	}
	return true
}

func validateObservedRange(or ObservedRange) {
	if model.NewPK(or.PK).IsZero() {
		panic("can't add observed range without pk")
	}
	if model.NewPK(or.LeaseNodePK).IsZero() {
		panic("can't add observed range without lease node id")
	}
	if model.NewPK(or.LeaseReplicaPK).IsZero() {
		panic("can't add observed range without lease replica id")
	}
	if reflect.ValueOf(or.Status).IsZero() {
		panic("can't add observed range without a status")
	}
}

func RetrieveAddOpenRanges(ctx context.Context, qExec query.Execute, o Observe) error {
	var openR []*models.Range
	if err := retrieveOpenRangesQuery(qExec, openR).Exec(ctx); err != nil {
		return err
	}
	for _, r := range openR {
		o.Add(ObservedRange{
			PK:             r.ID,
			LeaseReplicaPK: r.RangeLease.RangeReplicaID,
			LeaseNodePK:    r.RangeLease.RangeReplica.NodeID,
			Status:         r.Status,
		})
	}
	return nil
}
