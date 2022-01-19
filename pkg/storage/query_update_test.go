package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryUpdate", func() {
	BeforeEach(createDummyModel)
	AfterEach(deleteDummyModel)
	Describe("Update an item", func() {
		var newDummyModel *storage.ChannelConfig
		BeforeEach(func() {
			newDummyModel = &storage.ChannelConfig{
				ID:   dummyModel.ID,
				Name: "New Name",
			}
		})
		It("Should update it without error", func() {
			err := dummyStorage.NewUpdate().Model(newDummyModel).WherePK(dummyModel.ID).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after update", func() {
			if err := dummyStorage.NewUpdate().Model(newDummyModel).WherePK(dummyModel.ID).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}
			m := &storage.ChannelConfig{}
			if err := dummyStorage.NewRetrieve().Model(m).WherePK(dummyModel.ID).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}
			Expect(m.ID).To(Equal(dummyModel.ID))
			Expect(m.Name).To(Equal(newDummyModel.Name))
		})
	})
})
