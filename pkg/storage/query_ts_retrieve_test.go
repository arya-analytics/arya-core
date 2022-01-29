package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryTsRetrieve", func() {
	Describe("Standard usage", func() {
		Describe("Retrieving a sample", func() {
			It("Should create the index if it doesn't exist", func() {
				sample := &storage.ChannelSample{}
				err := mockStorage.NewTSRetrieve().Model(sample).WherePK(
					mockChannelCfg.ID).Exec(mockCtx)
				Expect(err).To(BeNil())
			})
		})
	})
})
