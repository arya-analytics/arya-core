package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
	BeforeEach(createDummyModel)
	Describe("Delete an item", func() {
		It("Should delete without error", func() {
			err := dummyEngine.NewDelete(dummyAdapter).Model(dummyModel).WherePK(
				dummyModel.
					ID).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should not be able to be re-queried after delete", func() {
			if err := dummyEngine.NewDelete(dummyAdapter).Model(dummyModel).WherePK(
				dummyModel.ID).Exec(dummyCtx); err != nil {
				log.Fatalln(err)
			}
			err := dummyEngine.NewRetrieve(dummyAdapter).Model(dummyModel).
				WherePK(dummyModel.ID).Exec(dummyCtx)
			Expect(err).ToNot(BeNil())
			Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
		})
	})
	Describe("Delete multiple items", func() {
		var err error
		var dummyModelTwo *storage.ChannelConfig
		AfterEach(func() {
			dummyEngine.NewDelete(dummyAdapter).Model(dummyModelTwo).WherePK(
				dummyModelTwo.ID).Exec(dummyCtx)
		})
		BeforeEach(func() {
			dummyModelTwo = &storage.ChannelConfig{
				ID:   uuid.New(),
				Name: "CC 45",
			}
			if err := dummyEngine.NewCreate(dummyAdapter).Model(dummyModelTwo).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}
			models := []*storage.ChannelConfig{dummyModel, dummyModelTwo}
			pks := []uuid.UUID{dummyModel.ID, dummyModelTwo.ID}
			err = dummyEngine.NewDelete(dummyAdapter).Model(&models).WherePKs(pks).
				Exec(dummyCtx)
		})
		It("Should delete them without error", func() {
			Expect(err).To(BeNil())
		})
		It("Should not be able to be-requeried after delete", func() {
			var models []*storage.ChannelConfig
			e := dummyEngine.NewRetrieve(dummyAdapter).Model(&models).WherePKs(
				[]uuid.UUID{dummyModelTwo.ID,
					dummyModel.ID}).Exec(dummyCtx)
			Expect(e).To(BeNil())
			Expect(models).To(HaveLen(0))
		})
	})
})
