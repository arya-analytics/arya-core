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
		It("Should delete without errutil", func() {
			err := mockEngine.NewDelete(mockAdapter).Model(mockModel).WherePK(
				mockModel.
					ID).Exec(
				mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should not be able to be re-queried after delete", func() {
			if err := mockEngine.NewDelete(mockAdapter).Model(mockModel).WherePK(
				mockModel.ID).Exec(mockCtx); err != nil {
				log.Fatalln(err)
			}
			err := mockEngine.NewRetrieve(mockAdapter).Model(mockModel).
				WherePK(mockModel.ID).Exec(mockCtx)
			Expect(err).ToNot(BeNil())
			Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
		})
	})
	Describe("Delete multiple items", func() {
		var err error
		var mockModelTwo *storage.ChannelConfig
		BeforeEach(func() {
			mockModelTwo = &storage.ChannelConfig{
				ID:     uuid.New(),
				Name:   "CC 45",
				NodeID: 1,
			}
			if err := mockEngine.NewCreate(mockAdapter).Model(mockModelTwo).Exec(
				mockCtx); err != nil {
				log.Fatalln(err)
			}
			models := []*storage.ChannelConfig{mockModel, mockModelTwo}
			pks := []uuid.UUID{mockModel.ID, mockModelTwo.ID}
			err = mockEngine.NewDelete(mockAdapter).Model(&models).WherePKs(pks).
				Exec(mockCtx)
		})
		It("Should delete them without errutil", func() {
			Expect(err).To(BeNil())
		})
		It("Should not be able to re-queried after delete", func() {
			var models []*storage.ChannelConfig
			e := mockEngine.NewRetrieve(mockAdapter).Model(&models).WherePKs(
				[]uuid.UUID{mockModelTwo.ID,
					mockModel.ID}).Exec(mockCtx)
			Expect(e).To(BeNil())
			Expect(models).To(HaveLen(0))
		})
	})
})
