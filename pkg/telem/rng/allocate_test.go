package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/telem/rng/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Allocate", func() {
	var (
		obs rng.Observe
		p   *mock.Persist
		svc *rng.Service
	)
	BeforeEach(func() {
		obs = rng.NewObserveMem([]rng.ObservedRange{})
		p = &mock.Persist{}
		svc = rng.NewService(obs, p)
	})
	Describe("A Chunk", func() {
		Context("When no open range is under observation", func() {
			It("Should allocate a new range", func() {
				chunkToAlloc := &models.ChannelChunk{}
				err := svc.NewAllocate().Chunk(1, chunkToAlloc).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(p.Ranges).To(HaveLen(1))
				Expect(obs.RetrieveAll()).To(HaveLen(1))
				_, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodeID: 1})
				Expect(ok).To(BeTrue())
			})
		})
		Context("When an open range is under observation", func() {
			It("Should allocate to the open range", func() {
				obs.Add(rng.ObservedRange{
					ID:             uuid.New(),
					Status:         models.RangeStatusOpen,
					LeaseNodeID:    1,
					LeaseReplicaID: uuid.New(),
				})
				chunkToAlloc := &models.ChannelChunk{}
				err := svc.NewAllocate().Chunk(1, chunkToAlloc).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(p.Ranges).To(HaveLen(0))
				Expect(obs.RetrieveAll()).To(HaveLen(1))
				or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodeID: 1})
				Expect(ok).To(BeTrue())
				Expect(chunkToAlloc.RangeID).To(Equal(or.ID))
			})
		})
	})
	Describe("A Chunk Replica", func() {
		Context("When a chunk hasn't been allocated", func() {
			It("Should panic", func() {
				chunkReplicaToAlloc := &models.ChannelChunkReplica{}
				Expect(func() {
					svc.NewAllocate().ChunkReplica(chunkReplicaToAlloc).Exec(ctx)
				}).To(Panic())
			})
		})
		Context("When a chunk has already been allocated", func() {
			Context("When the range remains open", func() {
				It("Should allocate to the open range", func() {
					chunkToAlloc := &models.ChannelChunk{}
					alloc := svc.NewAllocate()
					err := alloc.Chunk(1, chunkToAlloc).Exec(ctx)
					Expect(err).To(BeNil())
					chunkReplicaToAlloc := &models.ChannelChunkReplica{}
					crErr := alloc.ChunkReplica(chunkReplicaToAlloc).Exec(ctx)
					Expect(crErr).To(BeNil())
					Expect(p.Ranges).To(HaveLen(1))
					Expect(obs.RetrieveAll()).To(HaveLen(1))
					or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodeID: 1})
					Expect(ok).To(BeTrue())
					Expect(chunkReplicaToAlloc.RangeReplicaID).To(Equal(or.LeaseReplicaID))
				})
			})
			Context("When the range has been closed in between allocating the chunk and the replica", func() {
				It("Should allocate a new range", func() {
					chunkToAlloc := &models.ChannelChunk{}
					alloc := svc.NewAllocate()
					err := alloc.Chunk(1, chunkToAlloc).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(p.Ranges).To(HaveLen(1))
					obs.Add(rng.ObservedRange{
						ID:             p.Ranges[0].ID,
						LeaseReplicaID: p.Ranges[0].RangeLease.RangeReplica.ID,
						LeaseNodeID:    p.Ranges[0].RangeLease.RangeReplica.NodeID,
						Status:         models.RangeStatusClosed,
					})
					chunkReplicaToAlloc := &models.ChannelChunkReplica{}
					crErr := alloc.ChunkReplica(chunkReplicaToAlloc).Exec(ctx)
					Expect(crErr).To(BeNil())
					Expect(p.Ranges).To(HaveLen(2))
					Expect(obs.RetrieveAll()).To(HaveLen(2))
					or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodeID: 1})
					Expect(ok).To(BeTrue())
					Expect(chunkReplicaToAlloc.RangeReplicaID).To(Equal(or.LeaseReplicaID))
				})
			})
		})
	})

})