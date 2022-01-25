package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryRetrieve", func() {
	Describe("Normal Operation", func() {
		BeforeEach(createMockModel)
		AfterEach(deleteMockModel)
		Describe("Retrieve an item", func() {
			var retrievedModel = &storage.ChannelChunk{}
			var err error
			BeforeEach(func() {
				err = mockEngine.NewRetrieve(mockAdapter).Model(retrievedModel).WherePK(mockModel.
					ID).Exec(mockCtx)
			})
			It("Should retrieve it without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				Expect(retrievedModel.Data).ToNot(BeNil())
				b := make([]byte, retrievedModel.Data.Size())
				_, err = retrievedModel.Data.Read(b)
				Expect(err.Error()).To(Equal("EOF"))
				Expect(b).To(Equal(mockBytes))
			})
		})
		Describe("Retrieve multiple items", func() {
			It("Should retrieve the correct items", func() {
				mockModelTwo := &storage.ChannelChunk{
					ID:   uuid.New(),
					Data: mock.NewObject([]byte("model two")),
				}
				if err := mockEngine.NewCreate(mockAdapter).Model(mockModelTwo).Exec(
					mockCtx); err != nil {
					log.Fatalln(err)
				}

				var models []*storage.ChannelChunk
				err := mockEngine.NewRetrieve(mockAdapter).Model(&models).WherePKs([]uuid.
					UUID{mockModel.ID, mockModelTwo.ID}).Exec(mockCtx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect(models[0].ID == mockModelTwo.ID || models[1].ID == mockModelTwo.ID).
					To(BeTrue())
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesnt exist", func() {
			It("Should return the correct error type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &storage.ChannelChunk{}
				err := mockEngine.NewRetrieve(mockAdapter).Model(m).WherePK(
					somePKThatDoesntExist).Exec(mockCtx)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
