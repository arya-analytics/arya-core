package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Describe("Create a New ChannelConfig", func() {
		AfterEach(deleteDummyModel)
		It("Should create it without error", func() {
			err := dummyStorage.NewCreate().Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			err := dummyStorage.NewCreate().Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
			m := &storage.ChannelConfig{}
			err = dummyStorage.NewRetrieve().Model(m).WhereID(dummyModel.ID).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.ID).To(Equal(dummyModel.ID))
		})
	})
})
