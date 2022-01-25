package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
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
