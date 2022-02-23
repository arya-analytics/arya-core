package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/telem/rng/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Partition", func() {
	Context("Over allocated range", func() {
		var (
			p                rng.Persist
			rngId            uuid.UUID
			part             *rng.Partition
			newRanges        []*models.Range
			sourceChunkCount int
		)
		BeforeEach(func() {
			p, rngId = mock.NewPersistOverallocatedRange()
			part = &rng.Partition{RangeID: rngId, Persist: p}
			chunks, err := p.RetrieveRangeChunks(ctx, rngId)
			Expect(err).To(BeNil())
			sourceChunkCount = len(chunks)
			newRanges, err = part.Exec(ctx)
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
	})
})
