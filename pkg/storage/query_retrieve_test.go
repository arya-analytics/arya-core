package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("retrieveQuery", func() {
	BeforeEach(createDummyModel)
	AfterEach(deleteDummyModel)
	Describe("Retrieve a channel config", func() {
		It("Should retrieve without error", func() {
			m := &storage.ChannelConfig{}
			err := dummyStorage.NewRetrieve().Model(m).WhereID(dummyModel.ID).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should retrieve the correct item", func() {
			m := &storage.ChannelConfig{}
			err := dummyStorage.NewRetrieve().Model(m).WhereID(dummyModel.ID).Exec(dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.ID).To(Equal(dummyModel.ID))
			Expect(m.Name).To(Equal(dummyModel.Name))
		})
	})
})
