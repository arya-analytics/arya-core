package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	BeforeEach(createMockSeries)
	AfterEach(deleteMockChannelCfg)
	Describe("Standard usage", func() {
		BeforeEach(func() { createMockSamples(4) })
		Describe("Retrieving a sample", func() {
			var sample = &storage.ChannelSample{}
			var err error
			BeforeEach(func() {
				err = mockStorage.NewTSRetrieve().Model(sample).WherePK(
					mockChannelCfg.ID).Exec(mockCtx)
			})
			It("Should retrieve the sample without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct sample", func() {
				Expect(sample.ChannelConfigID).To(Equal(mockChannelCfg.ID))
			})
		})
		Describe("Retrieving a sample by time range", func() {
			var samples []*storage.ChannelSample
			var err error
			var tStart, tEnd time.Time
			BeforeEach(func() {
				sampleTime := time.Unix(0, mockSamples[1].Timestamp)
				tStart = sampleTime.Add(-1100 * time.Millisecond)
				tEnd = sampleTime.Add(500 * time.Millisecond)
				err = mockStorage.NewTSRetrieve().Model(&samples).WherePK(
					mockChannelCfg.ID).WhereTimeRange(tStart.UnixNano(), tEnd.UnixNano()).
					Exec(mockCtx)
			})
			It("Should retrieve the sample without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct samples", func() {
				Expect(samples).To(HaveLen(2))
				for _, sample := range samples {
					Expect(sample.Timestamp).Should(BeNumerically("<", tEnd.UnixNano()))
					Expect(sample.Timestamp).Should(BeNumerically(">", tStart.UnixNano()))
				}
			})
		})
	})
})
