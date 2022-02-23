package rng

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"reflect"
	"sync"
)

type Observe interface {
	Add(or ObservedRange)
	Retrieve(q ObservedRange) (ObservedRange, bool)
	RetrieveAll() []ObservedRange
	RetrieveFilter(q ObservedRange) []ObservedRange
}

type ObservedRange struct {
	ID             uuid.UUID
	Status         models.RangeStatus
	LeaseNodeID    int
	LeaseReplicaID uuid.UUID
}

type ObserveMem struct {
	mu     sync.Mutex
	ranges map[uuid.UUID]ObservedRange
}

func NewObserveMem(ranges []ObservedRange) *ObserveMem {
	rangeMap := map[uuid.UUID]ObservedRange{}
	for _, rng := range ranges {
		validateObservedRange(rng)
		rangeMap[rng.ID] = rng
	}
	return &ObserveMem{ranges: rangeMap}
}

func (o *ObserveMem) Add(or ObservedRange) {
	o.mu.Lock()
	defer o.mu.Unlock()
	validateObservedRange(or)
	o.ranges[or.ID] = or
}

func (o *ObserveMem) Retrieve(q ObservedRange) (ObservedRange, bool) {
	matches := o.RetrieveFilter(q)
	if len(matches) > 0 {
		return matches[0], true
	}
	return ObservedRange{}, false
}

func (o *ObserveMem) RetrieveAll() (ranges []ObservedRange) {
	// Copying here to protect against modification
	for _, v := range o.ranges {
		ranges = append(ranges, v)
	}
	return ranges
}

func (o *ObserveMem) RetrieveFilter(q ObservedRange) (matches []ObservedRange) {
	for _, or := range o.ranges {
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
	if !model.NewPK(q.ID).IsZero() && q.ID != or.ID {
		return false
	}
	if !model.NewPK(q.LeaseReplicaID).IsZero() && q.LeaseReplicaID != or.LeaseReplicaID {
		return false
	}
	if !model.NewPK(q.LeaseNodeID).IsZero() && q.LeaseNodeID != or.LeaseNodeID {
		return false
	}
	return true
}

func validateObservedRange(or ObservedRange) {
	if model.NewPK(or.ID).IsZero() {
		panic("can't add observed range without pk")
	}
	if model.NewPK(or.LeaseNodeID).IsZero() {
		panic("can't add observed range without lease node id")
	}
	if model.NewPK(or.LeaseReplicaID).IsZero() {
		panic("can't add observed range without lease replica id")
	}
	if reflect.ValueOf(or.Status).IsZero() {
		panic("can't add observed range without a status")
	}
}
