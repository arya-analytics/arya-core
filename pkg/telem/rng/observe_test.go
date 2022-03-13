package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("observe", func() {
	Describe("ObserveMem", func() {
		It("Should validate the observed rngMap without panicking", func() {
			Expect(func() {
				rng.NewObserveMem([]rng.ObservedRange{
					{PK: uuid.New(),
						Status:         models.RangeStatusClosed,
						LeaseNodePK:    1,
						LeaseReplicaPK: uuid.New()},
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
			Entry("No PK", rng.ObservedRange{LeaseReplicaPK: uuid.New(), LeaseNodePK: 1, Status: models.RangeStatusClosed}),
			Entry("No Lease Node PK", rng.ObservedRange{PK: uuid.New(), LeaseReplicaPK: uuid.New(), Status: models.RangeStatusClosed}),
			Entry("No Lease Replica PK", rng.ObservedRange{PK: uuid.New(), LeaseNodePK: 1, Status: models.RangeStatusOpen}),
			Entry("No Range Status", rng.ObservedRange{PK: uuid.New(), LeaseReplicaPK: uuid.New(), LeaseNodePK: 1}),
		)
		Describe("Retrieving rngMap", func() {
			var (
				or = rng.ObservedRange{
					PK:             uuid.New(),
					Status:         models.RangeStatusOpen,
					LeaseNodePK:    1,
					LeaseReplicaPK: uuid.New(),
				}
				ranges = []rng.ObservedRange{
					{
						PK:             uuid.New(),
						Status:         models.RangeStatusClosed,
						LeaseNodePK:    2,
						LeaseReplicaPK: uuid.New(),
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
						Expect(retOR.PK).To(Equal(or.PK))
					},
					Entry("By PK", ranges, or, rng.ObservedRange{PK: or.PK}),
					Entry("By Status", ranges, or, rng.ObservedRange{Status: or.Status}),
					Entry("By Lease Node PK", ranges, or, rng.ObservedRange{LeaseNodePK: or.LeaseNodePK}),
					Entry("by Lease Node PK and status", ranges, or, rng.ObservedRange{LeaseNodePK: or.LeaseNodePK, Status: models.RangeStatusOpen}),
					Entry("By Lease Replica PK", ranges, or, rng.ObservedRange{LeaseReplicaPK: or.LeaseReplicaPK}),
				)
				It("Should return false when a range can't be found", func() {
					obs := rng.NewObserveMem(ranges)
					_, ok := obs.Retrieve(rng.ObservedRange{PK: uuid.New()})
					Expect(ok).To(BeFalse())
				})
			})
			Describe("Retrieving all rngMap", func() {
				It("Should retrieve all rngMap correctly", func() {
					or := rng.NewObserveMem(ranges)
					Expect(or.RetrieveAll()).To(HaveLen(2))
				})
			})
		})
		Describe("It shouldn't run into any race conditions with frequent updates to the same item from diff goroutines", func() {
			wg := sync.WaitGroup{}
			or := rng.ObservedRange{
				PK:             uuid.New(),
				LeaseNodePK:    1,
				LeaseReplicaPK: uuid.New(),
				Status:         models.RangeStatusOpen,
			}
			obs := rng.NewObserveMem([]rng.ObservedRange{or})
			wg.Add(100)
			for i := 0; i < 100; i++ {
				go func(or rng.ObservedRange) {
					or.LeaseReplicaPK = uuid.New()
					obs.Add(or)
					wg.Done()
				}(or)
			}
		})
	})
})
