package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
	BeforeEach(createMockModel)
	Describe("Delete an item", func() {
		It("Should delete without error", func() {
			err := mockEngine.NewDelete(mockADapter).Model(mockModel).WherePK(
				mockModel.
					ID).Exec(
				mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should not be able to be re-queried after delete", func() {
			if err := mockEngine.NewDelete(mockADapter).Model(mockModel).WherePK(
				mockModel.ID).Exec(mockCtx); err != nil {
				log.Fatalln(err)
			}
			err := mockEngine.NewRetrieve(mockADapter).Model(mockModel).
				WherePK(mockModel.ID).Exec(mockCtx)
			Expect(err).ToNot(BeNil())
			Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
		})
	})
	Describe("Delete multiple items", func() {
		var err error
		var dummyModelTwo *storage.ChannelConfig
		AfterEach(func() {
			mockEngine.NewDelete(mockADapter).Model(dummyModelTwo).WherePK(
				dummyModelTwo.ID).Exec(mockCtx)
		})
		BeforeEach(func() {
			dummyModelTwo = &storage.ChannelConfig{
				ID:     uuid.New(),
				Name:   "CC 45",
				NodeID: 1,
			}
			if err := mockEngine.NewCreate(mockADapter).Model(dummyModelTwo).Exec(
				mockCtx); err != nil {
				log.Fatalln(err)
			}
			models := []*storage.ChannelConfig{mockModel, dummyModelTwo}
			pks := []uuid.UUID{mockModel.ID, dummyModelTwo.ID}
			err = mockEngine.NewDelete(mockADapter).Model(&models).WherePKs(pks).
				Exec(mockCtx)
		})
		It("Should delete them without error", func() {
			Expect(err).To(BeNil())
		})
		It("Should not be able to be-requeried after delete", func() {
			var models []*storage.ChannelConfig
			e := mockEngine.NewRetrieve(mockADapter).Model(&models).WherePKs(
				[]uuid.UUID{dummyModelTwo.ID,
					mockModel.ID}).Exec(mockCtx)
			Expect(e).To(BeNil())
			Expect(models).To(HaveLen(0))
		})
	})
})
