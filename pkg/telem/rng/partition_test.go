package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/telem/rng/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Partition", func() {
	Describe("PartitionExecute", func() {
		Context("Over Allocated Range", func() {
			var (
				p                         rng.Persist
				rngId                     uuid.UUID
				part                      *rng.PartitionExecute
				newRanges                 []*models.Range
				sourceChunkReplicaNodeIDs map[uuid.UUID]int
				sourceChunkCount          int
				sourceChunkReplicaCount   int
			)
			BeforeEach(func() {
				sourceChunkReplicaNodeIDs = map[uuid.UUID]int{}
				p, rngId = mock.NewPersistOverallocatedRange()
				part = rng.NewPartitionExecute(ctx, clust.Exec, rngId)
				chunks, err := p.RetrieveRangeChunks(ctx, rngId)
				Expect(err).To(BeNil())
				rangeReplicas, err := p.RetrieveRangeReplicas(ctx, rngId)
				Expect(err).To(BeNil())
				sourceChunkCount = len(chunks)
				chunkReplicas, err := p.RetrieveRangeChunkReplicas(ctx, rngId)
				Expect(len(chunkReplicas)).To(BeNumerically(">", 0))
				for _, ccr := range chunkReplicas {
					rr, ok := findRangeReplica(ccr.RangeReplicaID, rangeReplicas)
					Expect(ok).To(BeTrue())
					sourceChunkReplicaNodeIDs[ccr.ID] = rr.NodeID
				}
				Expect(err).To(BeNil())
				sourceChunkReplicaCount = len(chunkReplicas)
				newRanges, err = part.Exec()
				Expect(err).To(BeNil())

			})
			Context("New Range Basic Checks", func() {
				It("Should create one new range", func() {
					Expect(newRanges).To(HaveLen(1))
				})
				Specify("Defined range, range lease, and lease replica", func() {
					newRng := newRanges[0]
					Expect(model.NewPK(newRng.ID).IsZero()).To(BeFalse())
					Expect(model.NewPK(newRng.RangeLease.ID).IsZero()).To(BeFalse())
					Expect(model.NewPK(newRng.RangeLease.RangeReplica.ID).IsZero()).To(BeFalse())
				})
				Specify("Lease on correct node", func() {
					newRng := newRanges[0]
					sourceRng, err := p.RetrieveRange(ctx, rngId)
					Expect(err).To(BeNil())
					Expect(newRng.RangeLease.RangeReplica.NodeID).To(Equal(sourceRng.RangeLease.RangeReplica.NodeID))
				})
			})
			Context("New Range Size", func() {
				It("Should be smaller than the max range size", func() {
					newRng := newRanges[0]
					size, err := p.RetrieveRangeSize(ctx, newRng.ID)
					Expect(err).To(BeNil())
					Expect(size).To(BeNumerically("<", models.MaxRangeSize))
				})
				It("Should be roughly 1/4 the size of the max range", func() {
					newRng := newRanges[0]
					size, err := p.RetrieveRangeSize(ctx, newRng.ID)
					Expect(err).To(BeNil())
					Expect(size).To(BeNumerically(">", float64(models.MaxRangeSize)*0.2))
					Expect(size).To(BeNumerically("<", float64(models.MaxRangeSize)*0.3))
				})
			})
			Context("Source range size", func() {
				It("Should be smaller than the max range size", func() {
					size, err := p.RetrieveRangeSize(ctx, rngId)
					Expect(err).To(BeNil())
					Expect(size).To(BeNumerically("<", models.MaxRangeSize))
				})
				It("Should be pretty close to the max range size", func() {
					size, err := p.RetrieveRangeSize(ctx, rngId)
					Expect(err).To(BeNil())
					Expect(size).To(BeNumerically(">", float64(models.MaxRangeSize)*0.98))
				})
			})
			Context("New Range Replicas", func() {
				Specify("There should be one new replica per source replica", func() {
					sourceReplicas, err := p.RetrieveRangeReplicas(ctx, rngId)
					Expect(err).To(BeNil())
					Expect(sourceReplicas).To(HaveLen(3))
					newReplicas, err := p.RetrieveRangeReplicas(ctx, newRanges[0].ID)
					Expect(len(newReplicas)).To(Equal(len(sourceReplicas)))
				})
				Specify("The node PK of each new range replica must the same as the original replica", func() {
					sourceReplicas, err := p.RetrieveRangeReplicas(ctx, rngId)
					Expect(err).To(BeNil())
					newReplicas, err := p.RetrieveRangeReplicas(ctx, newRanges[0].ID)
					Expect(err).To(BeNil())
					for _, newRR := range newReplicas {
						match := false
						for _, sourceRR := range sourceReplicas {
							if newRR.NodeID == sourceRR.NodeID {
								match = true
							}
						}
						Expect(match).To(BeTrue())
					}
				})
			})
			Context("Reallocated chunks", func() {
				Specify("The amount of chunks in the source range and the new range should equal the total chunks", func() {
					sourceChunks, err := p.RetrieveRangeChunks(ctx, rngId)
					Expect(err).To(BeNil())
					newChunks, err := p.RetrieveRangeChunks(ctx, newRanges[0].ID)
					Expect(len(sourceChunks)).To(BeNumerically("<", sourceChunkCount))
					Expect(len(sourceChunks) + len(newChunks)).To(Equal(sourceChunkCount))
				})
			})
			Context("Reallocated ChunkData Replicas", func() {
				Specify("The amount of chunk replicas should remain the same", func() {
					sourceCCRs, err := p.RetrieveRangeChunkReplicas(ctx, rngId)
					Expect(err).To(BeNil())
					newCCRs, err := p.RetrieveRangeChunkReplicas(ctx, newRanges[0].ID)
					Expect(len(sourceCCRs)).To(BeNumerically("<", sourceChunkReplicaCount))
					Expect(len(sourceCCRs) + len(newCCRs)).To(Equal(sourceChunkReplicaCount))
				})
				Specify("Each reallocated chunk replica should belong to the same node as it did before realloc", func() {
					newReplicas, err := p.RetrieveRangeReplicas(ctx, newRanges[0].ID)
					Expect(err).To(BeNil())
					newCCRs, err := p.RetrieveRangeChunkReplicas(ctx, newRanges[0].ID)
					Expect(err).To(BeNil())
					for _, newCCR := range newCCRs {
						repl, ok := findRangeReplica(newCCR.RangeReplicaID, newReplicas)
						Expect(ok).To(BeTrue())
						sourceNodeID, ok := sourceChunkReplicaNodeIDs[newCCR.ID]
						Expect(ok).To(BeTrue())
						Expect(repl.NodeID).To(Equal(sourceNodeID))
					}
				})

			})
			Context("Range Status", func() {
				It("Should close the source range", func() {
					sourceRng, err := p.RetrieveRange(ctx, rngId)
					Expect(err).To(BeNil())
					Expect(sourceRng.Status).To(Equal(models.RangeStatusClosed))
				})
				It("Should open the new range", func() {
					newRng, err := p.RetrieveRange(ctx, newRanges[0].ID)
					Expect(err).To(BeNil())
					Expect(newRng.Status).To(Equal(models.RangeStatusOpen))
				})
			})
		})
	})
	Describe("PartitionScheduler", func() {
		const (
			accel        = 160
			tickInterval = 900 * time.Millisecond
		)
		var (
			rngPK uuid.UUID
			pst   rng.Persist
			svc   *rng.Service
			obs   rng.Observe
		)
		JustBeforeEach(func() {
			go svc.Start(ctx, tasks.ScheduleWithAccel(accel), tasks.ScheduleWithSilence())
			time.Sleep(10 * time.Millisecond)
			var err error
			go func() {
				err = <-svc.Errors()
			}()
			time.Sleep(tickInterval)
			svc.Stop()
			Expect(err).To(BeNil())
		})
		Describe("Persisted Detection", func() {
			BeforeEach(func() {
				pst, rngPK = mock.NewPersistOverallocatedRange()
				obs = rng.NewObserveMem([]rng.ObservedRange{})
				svc = rng.NewService(obs, pst)
			})
			It("Should reallocate the overallocated range", func() {
				size, err := pst.RetrieveRangeSize(ctx, rngPK)
				Expect(err).To(BeNil())
				Expect(size).To(BeNumerically("<", models.MaxRangeSize))
			})
			It("Should add the new range to the observe", func() {
				Expect(obs.RetrieveAll()).To(HaveLen(2))
				Expect(obs.RetrieveFilter(rng.ObservedRange{Status: models.RangeStatusOpen})).To(HaveLen(1))
			})
			It("Should close the source range", func() {
				Expect(obs.RetrieveFilter(rng.ObservedRange{Status: models.RangeStatusClosed})).To(HaveLen(1))
				sourceRng, ok := obs.Retrieve(rng.ObservedRange{PK: rngPK})
				Expect(ok).To(BeTrue())
				Expect(sourceRng.Status).To(Equal(models.RangeStatusClosed))
			})
		})
		Describe("Observed Detection", func() {
			BeforeEach(func() {
				pst, rngPK = mock.NewPersistOverallocatedRange()
				newRng, err := pst.RetrieveRange(ctx, rngPK)
				Expect(err).To(BeNil())
				obs = rng.NewObserveMem([]rng.ObservedRange{
					{
						PK:             newRng.ID,
						Status:         models.RangeStatusOpen,
						LeaseNodePK:    newRng.RangeLease.RangeReplica.NodeID,
						LeaseReplicaPK: newRng.RangeLease.RangeReplica.ID,
					},
				})
				svc = rng.NewService(obs, pst)
			})
			It("Should reallocate the overallocated range", func() {
				size, err := pst.RetrieveRangeSize(ctx, rngPK)
				Expect(err).To(BeNil())
				Expect(size).To(BeNumerically("<", models.MaxRangeSize))
			})
			It("Should only create one new range", func() {
				Expect(obs.RetrieveAll()).To(HaveLen(2))
			})
			It("Should add the new range to the observe", func() {
				Expect(obs.RetrieveAll()).To(HaveLen(2))
				Expect(obs.RetrieveFilter(rng.ObservedRange{Status: models.RangeStatusOpen})).To(HaveLen(1))
			})
			It("Should close the source range", func() {
				Expect(obs.RetrieveFilter(rng.ObservedRange{Status: models.RangeStatusClosed})).To(HaveLen(1))
				sourceRng, ok := obs.Retrieve(rng.ObservedRange{PK: rngPK})
				Expect(ok).To(BeTrue())
				Expect(sourceRng.Status).To(Equal(models.RangeStatusClosed))
			})
		})
	})
})

func findRangeReplica(rangeReplicaID uuid.UUID, rangeReplicas []*models.RangeReplica) (*models.RangeReplica, bool) {
	for _, rr := range rangeReplicas {
		if rr.ID == rangeReplicaID {
			return rr, true
		}
	}
	return nil, false
}
