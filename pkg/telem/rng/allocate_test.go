package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Allocate", func() {
	var (
		obs rng.Observe
		svc *rng.Service
		ds  *mock.DataSourceMem
	)
	BeforeEach(func() {
		ds = mock.NewDataSourceMem()
		obs = rng.NewObserveMem([]rng.ObservedRange{})
		svc = rng.NewService(obs, ds.Exec)
	})
	Describe("A ChunkData", func() {
		Context("When no open range is under observation", func() {
			It("Should Allocate a new range", func() {
				chunkToAlloc := &models.ChannelChunk{}
				err := svc.NewAllocate().Chunk(1, chunkToAlloc).Exec(ctx)
				Expect(err).To(BeNil())
				var resRng []*models.Range
				Expect(ds.NewRetrieve().Model(&resRng).Exec(ctx)).To(BeNil())
				Expect(resRng).To(HaveLen(1))
				Expect(obs.RetrieveAll()).To(HaveLen(1))
				_, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodePK: 1})
				Expect(ok).To(BeTrue())
			})
		})
		Context("When an open range is under observation", func() {
			It("Should Allocate to the open range", func() {
				obs.Add(rng.ObservedRange{
					PK:             uuid.New(),
					Status:         models.RangeStatusOpen,
					LeaseNodePK:    1,
					LeaseReplicaPK: uuid.New(),
				})
				chunkToAlloc := &models.ChannelChunk{}
				err := svc.NewAllocate().Chunk(1, chunkToAlloc).Exec(ctx)
				Expect(err).To(BeNil())
				rErr := ds.NewRetrieve().Model(&models.Range{}).Exec(ctx)
				Expect(rErr.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
				Expect(obs.RetrieveAll()).To(HaveLen(1))
				or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodePK: 1})
				Expect(ok).To(BeTrue())
				Expect(chunkToAlloc.RangeID).To(Equal(or.PK))
			})
		})
	})
	Describe("A ChunkData Replica", func() {
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
				It("Should Allocate to the open range", func() {
					chunkToAlloc := &models.ChannelChunk{}
					alloc := svc.NewAllocate()
					err := alloc.Chunk(1, chunkToAlloc).Exec(ctx)
					Expect(err).To(BeNil())
					chunkReplicaToAlloc := &models.ChannelChunkReplica{}
					crErr := alloc.ChunkReplica(chunkReplicaToAlloc).Exec(ctx)
					Expect(crErr).To(BeNil())
					var resRng []*models.Range
					Expect(ds.NewRetrieve().Model(&resRng).Exec(ctx)).To(BeNil())
					Expect(resRng).To(HaveLen(1))
					Expect(obs.RetrieveAll()).To(HaveLen(1))
					or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodePK: 1})
					Expect(ok).To(BeTrue())
					Expect(chunkReplicaToAlloc.RangeReplicaID).To(Equal(or.LeaseReplicaPK))
				})
			})
			Context("When the range has been closed in between allocating the chunk and the replica", func() {
				It("Should Allocate a new range", func() {
					chunkToAlloc := &models.ChannelChunk{}
					alloc := svc.NewAllocate()
					err := alloc.Chunk(1, chunkToAlloc).Exec(ctx)
					Expect(err).To(BeNil())
					resRng := &models.Range{}
					Expect(ds.NewRetrieve().Model(resRng).Exec(ctx))
					obs.Add(rng.ObservedRange{
						PK:             resRng.ID,
						LeaseReplicaPK: resRng.RangeLease.RangeReplica.ID,
						LeaseNodePK:    resRng.RangeLease.RangeReplica.NodeID,
						Status:         models.RangeStatusClosed,
					})
					chunkReplicaToAlloc := &models.ChannelChunkReplica{}
					crErr := alloc.ChunkReplica(chunkReplicaToAlloc).Exec(ctx)
					Expect(crErr).To(BeNil())
					var resRanges []*models.Range
					Expect(ds.NewRetrieve().Model(&resRanges).Exec(ctx)).To(BeNil())
					Expect(resRanges).To(HaveLen(2))
					Expect(obs.RetrieveAll()).To(HaveLen(2))
					or, ok := obs.Retrieve(rng.ObservedRange{Status: models.RangeStatusOpen, LeaseNodePK: 1})
					Expect(ok).To(BeTrue())
					Expect(chunkReplicaToAlloc.RangeReplicaID).To(Equal(or.LeaseReplicaPK))
				})
			})
		})
	})

})
