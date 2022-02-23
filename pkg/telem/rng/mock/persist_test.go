package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Persist", func() {
	var (
		p   *mock.Persist
		ctx context.Context
	)
	BeforeEach(func() {
		p = mock.NewBlankPersist()
		ctx = context.Background()
	})
	Describe("New Range", func() {
		It("Should create a new range, range replica, and lease", func() {
			rng, err := p.NewRange(ctx, 1)
			Expect(err).To(BeNil())
			Expect(rng.RangeLease.RangeReplica.NodeID).To(Equal(1))
		})
	})
	Describe("Create Range", func() {
		It("Should add the range to the the list of ranges", func() {
			err := p.CreateRange(ctx, &models.Range{ID: uuid.New()})
			Expect(err).To(BeNil())
			Expect(p.Ranges).To(HaveLen(1))
		})
	})
	Describe("Create Range Lease", func() {
		It("Should add the lease to the list of leases", func() {
			err := p.CreateRangeLease(ctx, &models.RangeLease{})
			Expect(err).To(BeNil())
			Expect(p.RangeLeases).To(HaveLen(1))
		})
	})
	Describe("Create Range Replica", func() {
		It("Should add the replica to the list of replicas", func() {
			err := p.CreateRangeReplica(ctx, &[]*models.RangeReplica{{}, {}, {}})
			Expect(err).To(BeNil())
			Expect(p.RangeReplicas).To(HaveLen(3))
		})
	})
	Describe("Retrieve Range", func() {
		It("Should retrieve the range correctly", func() {
			rng, err := p.NewRange(ctx, 1)
			Expect(err).To(BeNil())
			Expect(rng.RangeLease.RangeReplica.NodeID).To(Equal(1))
			rRng, err := p.RetrieveRange(ctx, rng.ID)
			Expect(rRng.ID).To(Equal(rng.ID))
			Expect(err).To(BeNil())
		})
	})
	Describe("Chunk Related Operations", func() {
		var (
			rngID uuid.UUID
			op    *mock.Persist
		)
		BeforeEach(func() {
			op, rngID = mock.NewPersistOverallocatedRange()
		})
		Describe("Retrieve Range Size", func() {
			It("Should retrieve the correct range size", func() {
				Expect(op.RetrieveRangeSize(ctx, rngID)).To(BeNumerically(">", models.MaxRangeSize))
				Expect(op.RetrieveRangeSize(ctx, rngID)).To(BeNumerically("<", float64(models.MaxRangeSize)*1.28))
			})
		})
		Describe("Retrieve Range Chunks", func() {
			It("Should retrieve the correct chunks", func() {
				chunks, err := op.RetrieveRangeChunks(ctx, rngID)
				Expect(err).To(BeNil())
				Expect(len(chunks)).To(BeNumerically(">", 0))
			})
		})
		Describe("Retrieve Range Chunk Replicas", func() {
			It("Should retrieve the correct replicas", func() {
				repls, err := op.RetrieveRangeChunkReplicas(ctx, rngID)
				Expect(err).To(BeNil())
				Expect(len(repls)).To(BeNumerically(">", 0))
			})
		})
		Describe("Retrieve Range Replicas", func() {
			It("Should retrieve the correct replicas", func() {
				repls, err := op.RetrieveRangeReplicas(ctx, rngID)
				Expect(err).To(BeNil())
				Expect(len(repls)).To(BeNumerically(">", 0))

			})
		})
		Describe("Reallocate Chunks", func() {
			It("Should reallocate the chunks correctly", func() {
				orChunks, err := op.RetrieveRangeChunks(ctx, rngID)
				Expect(err).To(BeNil())
				newRng, err := op.NewRange(ctx, 2)
				err = op.ReallocateChunks(ctx, model.NewReflect(&orChunks).PKChain().Raw(), newRng.ID)
				reChunks, err := op.RetrieveRangeChunks(ctx, newRng.ID)
				Expect(len(reChunks)).To(Equal(len(orChunks)))
				Expect(reChunks[0].RangeID).To(Equal(newRng.ID))
				size, err := op.RetrieveRangeSize(ctx, rngID)
				Expect(size).To(BeNumerically("<", models.MaxRangeSize))
			})
		})
		Describe("Reallocate Chunk Replicas", func() {
			It("Should reallocate chunk replicas correctly", func() {
				orChunks, err := op.RetrieveRangeChunks(ctx, rngID)
				orChunkReplicas, err := op.RetrieveRangeChunkReplicas(ctx, rngID)
				Expect(err).To(BeNil())
				newRng, err := op.NewRange(ctx, 2)
				err = op.ReallocateChunks(ctx, model.NewReflect(&orChunks).PKChain().Raw(), newRng.ID)
				Expect(err).To(BeNil())
				err = op.ReallocateChunkReplicas(ctx, model.NewReflect(&orChunkReplicas).PKChain().Raw(), newRng.RangeLease.RangeReplicaID)
				Expect(err).To(BeNil())
				reReplicas, err := op.RetrieveRangeChunkReplicas(ctx, newRng.ID)
				Expect(len(reReplicas)).To(Equal(len(orChunkReplicas)))

			})
		})
	})
})
