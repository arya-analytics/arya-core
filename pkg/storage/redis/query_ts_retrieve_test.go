package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	Describe("Standard Usage", func() {
		Describe("Retrieving the most recent sample", func() {
			var sample = &storage.ChannelSample{}
			var err error
			BeforeEach(func() {
				createMockSamples(1)
				err = mockEngine.NewTSRetrieve(mockAdapter).Model(sample).WherePK(
					mockSeries.ID).Exec(mockCtx)
			})
			It("Should retrieve without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				Expect(sample.ChannelConfigID).To(Equal(mockSeries.ID))
				Expect(sample.Timestamp).To(Equal(mockSamples[0].Timestamp))
				Expect(sample.Value).To(Equal(mockSamples[0].Value))
			})
		})
		Describe("Retrieving all samples", func() {
			var samples []*storage.ChannelSample
			var err error
			BeforeEach(func() {
				createMockSamples(3)
				err = mockEngine.NewTSRetrieve(mockAdapter).Model(&samples).WherePK(
					mockSeries.ID).AllTimeRange().Exec(mockCtx)
			})
			It("Should retrieve without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct items", func() {
				Expect(samples).To(HaveLen(3))
			})
		})
		Describe("Retrieving samples from multiple pks", func() {
			var err error
			var samples []*storage.ChannelSample
			var mockSampleTwo *storage.ChannelSample
			BeforeEach(func() {
				createMockSamples(3)
				mockSeriesTwo := &storage.ChannelConfig{
					Name: "SG_03",
					ID:   uuid.New(),
				}
				if cErrOne := mockEngine.NewTSCreate(mockAdapter).Series().Model(mockSeriesTwo).Exec(
					mockCtx); cErrOne != nil {
					panic(cErrOne)
				}
				mockSampleTwo = &storage.ChannelSample{
					Timestamp:       time.Now().UnixNano(),
					Value:           123.2,
					ChannelConfigID: mockSeriesTwo.ID,
				}
				if cErrTwo := mockEngine.NewTSCreate(mockAdapter).Sample().Model(mockSampleTwo).
					Exec(mockCtx); cErrTwo != nil {
					panic(cErrTwo)
				}
				err = mockEngine.NewTSRetrieve(mockAdapter).Model(&samples).WherePKs(
					[]uuid.UUID{mockSeriesTwo.ID, mockSeries.ID}).AllTimeRange().Exec(
					mockCtx)
			})
			It("Should retrieve without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct items", func() {
				Expect(samples).To(HaveLen(4))
			})
		})
		Describe("Retrieve samples across a time range", func() {
			var err error
			var samples []*storage.ChannelSample
			var mockSampleTwo *storage.ChannelSample
			BeforeEach(func() {
				createMockSamples(3)
				mockSeriesTwo := &storage.ChannelConfig{
					Name: "SG_03",
					ID:   uuid.New(),
				}
				if cErrOne := mockEngine.NewTSCreate(mockAdapter).Series().Model(mockSeriesTwo).Exec(
					mockCtx); cErrOne != nil {
					panic(cErrOne)
				}
				mockSampleTwo = &storage.ChannelSample{
					Timestamp:       0,
					Value:           123.2,
					ChannelConfigID: mockSeriesTwo.ID,
				}
				if cErrTwo := mockEngine.NewTSCreate(mockAdapter).Sample().Model(mockSampleTwo).
					Exec(mockCtx); cErrTwo != nil {
					panic(cErrTwo)
				}
				toTS := time.Now().Add(3 * time.Second).UnixNano()
				fromTS := time.Now().Add(-15 * time.Second).UnixNano()
				err = mockEngine.NewTSRetrieve(mockAdapter).Model(&samples).WherePKs(
					[]uuid.UUID{mockSeriesTwo.ID, mockSeries.ID}).WhereTimeRange(fromTS,
					toTS).
					Exec(mockCtx)
			})
			It("Should retrieve without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct items", func() {
				Expect(samples).To(HaveLen(3))
			})
		})
		Describe("Checking if a series exists", func() {
			Context("The series does not exist", func() {
				It("Should return false", func() {
					e, err := mockEngine.NewTSRetrieve(mockAdapter).SeriesExists(
						mockCtx, uuid.New())
					Expect(e).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})
		})
	})
	Describe("Edge cases + errors", func() {
		BeforeEach(func() { createMockSamples(1) })
		Context("Retrieving a sample", func() {
			s := &storage.ChannelSample{}
			Context("No PK provided", func() {
				It("Should return the correct storage error", func() {
					err := mockEngine.NewTSRetrieve(mockAdapter).Model(s).Exec(mockCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeInvalidArgs))
				})
			})
			Context("Invalid PK provided", func() {
				It("Should return the correct storage error", func() {
					err := mockEngine.NewTSRetrieve(mockAdapter).WherePK(uuid.New()).Model(s).
						Exec(mockCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
				})
			})
		})
	})
})
