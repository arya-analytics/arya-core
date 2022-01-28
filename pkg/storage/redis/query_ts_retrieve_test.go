package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryTsRetrieve", func() {
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
})
