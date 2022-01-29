package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryTsCreate", func() {
	BeforeEach(createMockChannelCfg)
	AfterEach(deleteMockChannelCfg)
	Describe("Standard Usage", func() {
		Describe("Create a new series", func() {
			// NOTE: If this entire Describe section is run,
			// err should end up being a unique violation.
			var err error
			BeforeEach(func() {
				err = mockStorage.NewTSCreate().Series().Model(mockChannelCfg).Exec(
					mockCtx)
			})
			It("Should create the series without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should exist after creation", func() {
				exists, rErr := mockStorage.NewTSRetrieve().SeriesExists(mockCtx,
					mockChannelCfg.ID)
				Expect(rErr).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
		Describe("Create a new sample", func() {

		})
	})
})
