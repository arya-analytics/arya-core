package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng/mock"
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
	Describe("New Range Replica", func() {
		It("Should create a new range replica", func() {
			rng, err := p.NewRange(ctx, 1)
			Expect(err).To(BeNil())
			rr, err := p.NewRangeReplica(ctx, rng.ID, 2)
			Expect(err).To(BeNil())
			Expect(rr.NodeID).To(Equal(2))

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
				Expect(repls[0].NodeID).To(BeNumerically(">", 0))

			})
		})
		Describe("Reallocate Chunks", func() {
			It("Should reallocate the chunks correctly", func() {
				orChunks, err := op.RetrieveRangeChunks(ctx, rngID)
				Expect(err).To(BeNil())
				newRng, err := op.NewRange(ctx, 2)
				var (
					ccPKs []uuid.UUID
				)
				for _, cc := range orChunks {
					ccPKs = append(ccPKs, cc.ID)
				}
				err = op.ReallocateChunks(ctx, ccPKs, newRng.ID)
				reChunks, err := op.RetrieveRangeChunks(ctx, newRng.ID)
				Expect(len(reChunks)).To(Equal(len(orChunks)))
				Expect(reChunks[0].RangeID).To(Equal(newRng.ID))
				size, err := op.RetrieveRangeSize(ctx, rngID)
				Expect(size).To(BeNumerically("<", models.MaxRangeSize))
				newSourceChunks, err := op.RetrieveRangeChunks(ctx, rngID)
				Expect(len(newSourceChunks)).To(Equal(0))
			})
		})
		Describe("Reallocate Chunk Replicas", func() {
			It("Should reallocate chunk replicas correctly", func() {
				orChunks, err := op.RetrieveRangeChunks(ctx, rngID)
				orChunkReplicas, err := op.RetrieveRangeChunkReplicas(ctx, rngID)
				Expect(err).To(BeNil())
				var (
					ccPKs  []uuid.UUID
					ccrPKs []uuid.UUID
				)
				for _, cc := range orChunks {
					ccPKs = append(ccPKs, cc.ID)
				}
				for _, ccr := range orChunks {
					ccrPKs = append(ccrPKs, ccr.ID)
				}
				newRng, err := op.NewRange(ctx, 2)
				err = op.ReallocateChunks(ctx, ccPKs, newRng.ID)
				Expect(err).To(BeNil())
				err = op.ReallocateChunkReplicas(ctx, ccrPKs, newRng.RangeLease.RangeReplicaID)
				Expect(err).To(BeNil())
				reReplicas, err := op.RetrieveRangeChunkReplicas(ctx, newRng.ID)
				Expect(len(reReplicas)).To(Equal(len(orChunkReplicas)))

			})
		})
	})
})
