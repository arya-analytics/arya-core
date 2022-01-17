package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
	BeforeEach(createDummyModel)
	Describe("Delete a channel config", func() {
		It("Should delete without error", func() {
			err := dummyStorage.NewDelete().Model(dummyModel).WherePK(dummyModel.ID).
				Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Shouldn't throw an error when trying to retrieve after deletion", func() {
			if err := dummyStorage.NewDelete().Model(dummyModel).WherePK(dummyModel.
				ID).Exec(dummyCtx); err != nil {
				log.Fatalln(err)
			}
			err := dummyStorage.NewRetrieve().Model(dummyModel).WherePK(dummyModel.
				ID).Exec(dummyCtx)
			Expect(err).ToNot(BeNil())
		})
	})
})
