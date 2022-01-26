package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
	Describe("Standard Usage", func() {
		BeforeEach(createMockChannelCfg)
		Describe("Delete a channel config", func() {
			It("Should delete without errutil", func() {
				err := mockStorage.NewDelete().Model(mockChannelCfg).WherePK(mockChannelCfg.ID).
					Exec(mockCtx)
				Expect(err).To(BeNil())
			})
			It("Shouldn't throw an errutil when trying to retrieve after deletion", func() {
				if err := mockStorage.NewDelete().Model(mockChannelCfg).WherePK(mockChannelCfg.
					ID).Exec(mockCtx); err != nil {
					log.Fatalln(err)
				}
				err := mockStorage.NewRetrieve().Model(mockChannelCfg).WherePK(mockChannelCfg.
					ID).Exec(mockCtx)
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Describe("Providing no where statement to the query", func() {
			It("Should return an error", func() {
				err := mockStorage.NewDelete().Model(&storage.ChannelConfig{}).Exec(mockCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeInvalidArgs))
			})
		})
	})
})
