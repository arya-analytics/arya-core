package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	})
})
