package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyModel = &storage.ChannelConfig{
	ID:   432,
	Name: "Cool Name",
}

var _ = Describe("Create", func() {
	Describe("Create a New ChannelConfig", func() {
		It("Should create it without error", func() {
			err := dummyStorage.NewCreate().Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			m := &storage.ChannelConfig{}
			err := dummyStorage.NewRetrieve().Model(m).WhereID(dummyModel.ID).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.ID).To(Equal(dummyModel.ID))
		})
	})
})
