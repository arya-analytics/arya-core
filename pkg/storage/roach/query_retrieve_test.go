package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryRetrieve", func() {
	Describe("Standard Usage", func() {
		BeforeEach(createMockModel)
		AfterEach(deleteMockModel)
		Describe("Retrieve an item", func() {
			It("Should retrieve it without errutil", func() {
				m := &storage.ChannelConfig{}
				err := mockEngine.NewRetrieve(mockAdapter).Model(m).WherePK(mockModel.
					ID).Exec(mockCtx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				m := &storage.ChannelConfig{}
				if err := mockEngine.NewRetrieve(mockAdapter).Model(m).WherePK(mockModel.
					ID).Exec(mockCtx); err != nil {
					log.Fatalln(err)
				}
				Expect(m.ID).To(Equal(mockModel.ID))
				Expect(m.Name).To(Equal(mockModel.Name))
			})
		})
		Describe("Retrieve multiple items", func() {
			It("Should retrieve all the correct items", func() {
				mockModelTwo := &storage.ChannelConfig{
					ID:     uuid.New(),
					Name:   "CC 45",
					NodeID: 1,
				}
				if err := mockEngine.NewCreate(mockAdapter).Model(mockModelTwo).Exec(
					mockCtx); err != nil {
					log.Fatalln(err)
				}

				var models []*storage.ChannelConfig
				err := mockEngine.NewRetrieve(mockAdapter).Model(&models).WherePKs(
					[]uuid.UUID{mockModelTwo.ID,
						mockModel.ID}).Exec(mockCtx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect(models[0].Name == mockModel.Name || models[0].
					Name == mockModelTwo.Name).To(BeTrue())
				Expect(models[1].ID == mockModelTwo.ID || models[1].ID == mockModel.
					ID).To(BeTrue())
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct errutil type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &storage.ChannelConfig{}
				err := mockEngine.NewRetrieve(mockAdapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(mockCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
			})
		})
	})
})
