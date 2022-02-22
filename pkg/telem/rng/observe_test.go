package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("Observe", func() {
	Describe("ObserveMem", func() {
		It("Should validate the observed ranges without panicking", func() {
			Expect(func() {
				rng.NewObserveMem([]rng.ObservedRange{
					{ID: uuid.New(),
						Status:         models.RangeStatusClosed,
						LeaseNodeID:    1,
						LeaseReplicaID: uuid.New()},
				})
			}).ToNot(Panic())
		})
		DescribeTable("Invalid observed range additions",
			func(or rng.ObservedRange) {
				obs := rng.NewObserveMem([]rng.ObservedRange{})
				Expect(func() {
					obs.Add(or)
				}).To(Panic())
			},
			Entry("No ID", rng.ObservedRange{LeaseReplicaID: uuid.New(), LeaseNodeID: 1, Status: models.RangeStatusClosed}),
			Entry("No Lease Node ID", rng.ObservedRange{ID: uuid.New(), LeaseReplicaID: uuid.New(), Status: models.RangeStatusClosed}),
			Entry("No Lease Replica ID", rng.ObservedRange{ID: uuid.New(), LeaseNodeID: 1, Status: models.RangeStatusOpen}),
			Entry("No Range Status", rng.ObservedRange{ID: uuid.New(), LeaseReplicaID: uuid.New(), LeaseNodeID: 1}),
		)
		Describe("Retrieving ranges", func() {
			var (
				or = rng.ObservedRange{
					ID:             uuid.New(),
					Status:         models.RangeStatusOpen,
					LeaseNodeID:    1,
					LeaseReplicaID: uuid.New(),
				}
				ranges = []rng.ObservedRange{
					{
						ID:             uuid.New(),
						Status:         models.RangeStatusClosed,
						LeaseNodeID:    2,
						LeaseReplicaID: uuid.New(),
					},
					or,
				}
			)
			Describe("Retrieving a single range", func() {
				DescribeTable("Different query patterns",
					func(ranges []rng.ObservedRange, or rng.ObservedRange, q rng.ObservedRange) {
						obs := rng.NewObserveMem(ranges)
						retOR, ok := obs.Retrieve(q)
						Expect(ok).To(BeTrue())
						Expect(retOR.ID).To(Equal(or.ID))
					},
					Entry("By ID", ranges, or, rng.ObservedRange{ID: or.ID}),
					Entry("By Status", ranges, or, rng.ObservedRange{Status: or.Status}),
					Entry("By Lease Node ID", ranges, or, rng.ObservedRange{LeaseNodeID: or.LeaseNodeID}),
					Entry("by Lease Node ID and status", ranges, or, rng.ObservedRange{LeaseNodeID: or.LeaseNodeID, Status: models.RangeStatusOpen}),
					Entry("By Lease Replica ID", ranges, or, rng.ObservedRange{LeaseReplicaID: or.LeaseReplicaID}),
				)
				It("Should return false when a range can't be found", func() {
					obs := rng.NewObserveMem(ranges)
					_, ok := obs.Retrieve(rng.ObservedRange{ID: uuid.New()})
					Expect(ok).To(BeFalse())
				})
			})
			Describe("Retrieving all ranges", func() {
				It("Should retrieve all ranges correctly", func() {
					or := rng.NewObserveMem(ranges)
					Expect(or.RetrieveAll()).To(HaveLen(2))
				})
			})
		})
		Describe("It shouldn't run into any race conditions with frequent updates to the same item from diff goroutines", func() {
			wg := sync.WaitGroup{}
			or := rng.ObservedRange{
				ID:             uuid.New(),
				LeaseNodeID:    1,
				LeaseReplicaID: uuid.New(),
				Status:         models.RangeStatusOpen,
			}
			obs := rng.NewObserveMem([]rng.ObservedRange{or})
			wg.Add(100)
			for i := 0; i < 100; i++ {
				go func(or rng.ObservedRange) {
					or.LeaseReplicaID = uuid.New()
					obs.Add(or)
					wg.Done()
				}(or)
			}
		})
	})
})
