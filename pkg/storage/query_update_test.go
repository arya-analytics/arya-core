package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryUpdate", func() {
	BeforeEach(createMockChannelCfg)
	AfterEach(deleteMockChannelCfg)
	Describe("Update an item", func() {
		var newDummyModel *storage.ChannelConfig
		BeforeEach(func() {
			newDummyModel = &storage.ChannelConfig{
				ID:     mockChannelCfg.ID,
				Name:   "New Name",
				NodeID: 1,
			}
		})
		It("Should update it without errutil", func() {
			err := mockStorage.NewUpdate().Model(newDummyModel).WherePK(mockChannelCfg.ID).Exec(
				mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after update", func() {
			if err := mockStorage.NewUpdate().Model(newDummyModel).WherePK(mockChannelCfg.ID).Exec(
				mockCtx); err != nil {
				log.Fatalln(err)
			}
			m := &storage.ChannelConfig{}
			if err := mockStorage.NewRetrieve().Model(m).WherePK(mockChannelCfg.ID).Exec(
				mockCtx); err != nil {
				log.Fatalln(err)
			}
			Expect(m.ID).To(Equal(mockChannelCfg.ID))
			Expect(m.Name).To(Equal(newDummyModel.Name))
		})
	})
})
