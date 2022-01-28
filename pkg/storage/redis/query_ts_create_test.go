package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsCreate", func() {
	Describe("Standard Usage", func() {
		BeforeEach(createMockSeries)
		Describe("Create a new series", func() {
			It("Should exist after creation", func() {
				exists, err := mockEngine.NewTSRetrieve(mockAdapter).SeriesExists(mockCtx,
					mockSeries.ID)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
		Describe("Create a new sample", func() {
			Context("Single sample", func() {
				It("Should be able to re-retrieve the sample after creation", func() {
					mockSample := &storage.ChannelSample{
						Timestamp:       time.Now().UnixNano(),
						Value:           123.2,
						ChannelConfigID: mockSeries.ID,
					}
					err := mockEngine.NewTSCreate(mockAdapter).Sample().Model(
						mockSample).Exec(mockCtx)
					Expect(err).To(BeNil())
					var samples []*storage.ChannelSample
					rErr := mockEngine.NewTSRetrieve(mockAdapter).Model(&samples).
						WherePK(mockSeries.ID).Exec(
						mockCtx)
					Expect(rErr).To(BeNil())
					Expect(samples).To(HaveLen(1))
					Expect(samples[0].ChannelConfigID).To(Equal(mockSeries.ID))
					Expect(samples[0].Value).To(Equal(mockSample.Value))
					Expect(samples[0].Timestamp).To(Equal(mockSample.Timestamp))
				})
			})
		})
	})
	Describe("Edge cases + error", func() {
		Describe("Not selecting a variant", func() {
			It("Should return the correct storage error", func() {
				err := mockEngine.NewTSCreate(mockAdapter).Model(mockSeries).Exec(mockCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeInvalidArgs))
			})
		})
	})
})
