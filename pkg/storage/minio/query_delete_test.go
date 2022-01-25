package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryDelete", func() {
	Describe("Normal Operation", func() {
		BeforeEach(createMockModel)
		Describe("Delete an item", func() {
			var err error
			BeforeEach(func() {
				err = mockEngine.NewDelete(mockAdapter).Model(mockModel).WherePK(
					mockModel.ID).Exec(mockCtx)
			})
			It("Should delete it without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should not be able to be re-queried after delete", func() {
				rErr := mockEngine.NewRetrieve(mockAdapter).Model(mockModel).WherePK(
					mockModel.ID).Exec(mockCtx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
			})
		})
		Describe("Delete multiple Items", func() {
			var err error
			var mockModelTwo *storage.ChannelChunk
			BeforeEach(func() {
				mockModelTwo = &storage.ChannelChunk{
					ID:   uuid.New(),
					Data: mock.NewObject([]byte("mock bytes")),
				}
				if err := mockEngine.NewCreate(mockAdapter).Model(mockModelTwo).Exec(
					mockCtx); err != nil {
					log.Fatalln(err)
				}
				pks := []uuid.UUID{mockModel.ID, mockModelTwo.ID}
				err = mockEngine.NewDelete(mockAdapter).Model(mockModelTwo).WherePKs(pks).
					Exec(mockCtx)
			})
			It("Should delete them without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should not be able to be re-queried after delete", func() {
				var models []*storage.ChannelChunk
				e := mockEngine.NewRetrieve(mockAdapter).Model(&models).WherePKs(
					[]uuid.UUID{mockModelTwo.ID,
						mockModel.ID}).Exec(mockCtx)
				Expect(e).To(BeNil())
				Expect(models).To(HaveLen(0))
			})

		})
	})
})
