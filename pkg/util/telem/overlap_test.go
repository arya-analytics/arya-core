package telem_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"

	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

var _ = Describe("Overlap", func() {
	Describe("ChunkOverlap", func() {
		Context("Valid", func() {
			Describe("Partial", func() {
				Context("Uniform", func() {
					var (
						cOne *telem.Chunk
						cTwo *telem.Chunk
						o    telem.ChunkOverlap
					)
					BeforeEach(func() {
						cdOne := telem.NewChunkData([]byte{})
						Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
						cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
						cdTwo := telem.NewChunkData([]byte{})
						Expect(cdTwo.WriteData([]float64{6, 7, 8, 9, 10, 11, 12})).To(BeNil())
						cTwo = telem.NewChunk(cOne.Start().Add(telem.NewTimeSpan(5*time.Second)), telem.DataTypeFloat64, telem.DataRate(1), cdTwo)
						Expect(cTwo.ValueAtTS(cTwo.Start())).To(Equal(float64(6)))
						o = cOne.Overlap(cTwo)
					})
					It("Should be valid", func() {
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should be uniform", func() {
						Expect(o.IsUniform()).To(BeTrue())
					})
					It("Should return the correct range", func() {
						Expect(o.Range().Start()).To(Equal(cTwo.Start()))
						Expect(o.Range().End()).To(Equal(cOne.End()))
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should return the correct source values", func() {
						Expect(o.SourceValues().([]float64)).To(Equal([]float64{6, 7, 8, 9}))
					})
					It("Should return the correct dest values", func() {
						Expect(o.DestValues().([]float64)).To(Equal([]float64{6, 7, 8, 9}))
					})
				})
				Context("Non-uniform", func() {
					var (
						cOne *telem.Chunk
						cTwo *telem.Chunk
						o    telem.ChunkOverlap
					)
					BeforeEach(func() {
						cdOne := telem.NewChunkData([]byte{})
						Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
						cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
						cdTwo := telem.NewChunkData([]byte{})
						Expect(cdTwo.WriteData([]float64{7, 8, 9, 10, 11, 12})).To(BeNil())
						cTwo = telem.NewChunk(cOne.Start().Add(telem.NewTimeSpan(5*time.Second)), telem.DataTypeFloat64, telem.DataRate(1), cdTwo)
						Expect(cTwo.ValueAtTS(cTwo.Start())).To(Equal(float64(7)))
						o = cOne.Overlap(cTwo)
					})
					It("Should be valid", func() {
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should be non-uniform", func() {
						Expect(o.IsUniform()).To(BeFalse())
					})
					It("Should return the correct range", func() {
						Expect(o.Range().Start()).To(Equal(cTwo.Start()))
						Expect(o.Range().End()).To(Equal(cOne.End()))
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should return the correct source values", func() {
						Expect(o.SourceValues().([]float64)).To(Equal([]float64{6, 7, 8, 9}))
					})
					It("Should return the correct dest values", func() {
						Expect(o.DestValues().([]float64)).To(Equal([]float64{7, 8, 9, 10}))
					})
				})
			})
			Describe("One chunk consumes another", func() {
				Context("Uniform", func() {
					var (
						cOne *telem.Chunk
						cTwo *telem.Chunk
						o    telem.ChunkOverlap
					)
					BeforeEach(func() {
						cdOne := telem.NewChunkData([]byte{})
						Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
						cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
						cdTwo := telem.NewChunkData([]byte{})
						Expect(cdTwo.WriteData([]float64{2, 3, 4, 5})).To(BeNil())
						cTwoStart := cOne.Start().Add(telem.NewTimeSpan(1 * time.Second))
						cTwo = telem.NewChunk(cTwoStart, telem.DataTypeFloat64, telem.DataRate(1), cdTwo)
						Expect(cTwo.ValueAtTS(cTwo.Start())).To(Equal(float64(2)))
						o = cOne.Overlap(cTwo)
					})
					It("Should be valid", func() {
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should be uniform", func() {
						Expect(o.IsUniform()).To(BeTrue())
					})
					It("Should return the correct values in the source range", func() {
						Expect(o.SourceValues()).To(Equal([]float64{2, 3, 4, 5}))
					})
					It("Should return the correct values in the dest range", func() {
						Expect(o.DestValues()).To(Equal([]float64{2, 3, 4, 5}))
					})
				})
				Context("Non-uniform", func() {
					var (
						cOne *telem.Chunk
						cTwo *telem.Chunk
						o    telem.ChunkOverlap
					)
					BeforeEach(func() {
						cdOne := telem.NewChunkData([]byte{})
						Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
						cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
						cdTwo := telem.NewChunkData([]byte{})
						Expect(cdTwo.WriteData([]float64{3, 4, 5, 6})).To(BeNil())
						cTwoStart := cOne.Start().Add(telem.NewTimeSpan(1 * time.Second))
						cTwo = telem.NewChunk(cTwoStart, telem.DataTypeFloat64, telem.DataRate(1), cdTwo)
						Expect(cTwo.ValueAtTS(cTwo.Start())).To(Equal(float64(3)))
						o = cOne.Overlap(cTwo)
					})
					It("Should be valid", func() {
						Expect(o.IsValid()).To(BeTrue())
					})
					It("Should be uniform", func() {
						Expect(o.IsUniform()).To(BeFalse())
					})
					It("Should return the correct values in the source range", func() {
						Expect(o.SourceValues()).To(Equal([]float64{2, 3, 4, 5}))
					})
					It("Should return the correct values in the dest range", func() {
						Expect(o.DestValues()).To(Equal([]float64{3, 4, 5, 6}))
					})
				})
			})
		})
		Context("InValid", func() {
			Context("No Overlap", func() {
				var (
					cOne *telem.Chunk
					cTwo *telem.Chunk
					o    telem.ChunkOverlap
				)
				BeforeEach(func() {
					cdOne := telem.NewChunkData([]byte{})
					Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
					cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
					cdTwo := telem.NewChunkData([]byte{})
					Expect(cdTwo.WriteData([]float64{11, 12})).To(BeNil())
					cTwoStart := cOne.Start().Add(telem.NewTimeSpan(11 * time.Second))
					cTwo = telem.NewChunk(cTwoStart, telem.DataTypeFloat64, telem.DataRate(1), cdTwo)
					o = cOne.Overlap(cTwo)
				})
				It("Should not be valid", func() {
					Expect(o.IsValid()).To(BeFalse())
				})
				It("Should be non-uniform", func() {
					Expect(o.IsUniform()).To(BeFalse())
				})
			})
			Describe("Incompatible data rates", func() {
				var (
					cOne *telem.Chunk
					cTwo *telem.Chunk
					o    telem.ChunkOverlap
				)
				BeforeEach(func() {
					cdOne := telem.NewChunkData([]byte{})
					Expect(cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
					cOne = telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cdOne)
					cdTwo := telem.NewChunkData([]byte{})
					Expect(cdTwo.WriteData([]float32{11, 12})).To(BeNil())
					cTwoStart := cOne.Start().Add(telem.NewTimeSpan(11 * time.Second))
					cTwo = telem.NewChunk(cTwoStart, telem.DataTypeFloat32, telem.DataRate(1), cdTwo)
					o = cOne.Overlap(cTwo)
				})
				It("Should not be valid", func() {
					Expect(o.IsValid()).To(BeFalse())
				})
				It("Should be non-uniform", func() {
					Expect(o.IsUniform()).To(BeFalse())
				})
			})
		})
		Describe("Modification", func() {
			Describe("Removing the ChunkOverlap", func() {
				var (
					cOne *telem.Chunk
					cTwo *telem.Chunk
				)
				BeforeEach(func() {
					cdOne := telem.NewChunkData([]byte{})
					cdOne.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})

					cOne = telem.NewChunk(
						telem.TimeStamp(0),
						telem.DataTypeFloat64,
						telem.DataRate(1),
						cdOne,
					)

					cdTwo := telem.NewChunkData([]byte{})
					cdTwo.WriteData([]float64{6, 7, 8, 9, 10, 11, 12})

					cTwo = telem.NewChunk(
						cOne.Start().Add(telem.NewTimeSpan(5*time.Second)),
						telem.DataTypeFloat64,
						telem.DataRate(1),
						cdTwo,
					)
					Expect(cTwo.ValueAtTS(cTwo.Start())).To(Equal(float64(6)))
				})
				It("Should remove the overlap correctly", func() {
				})
			})
		})
	})
})
