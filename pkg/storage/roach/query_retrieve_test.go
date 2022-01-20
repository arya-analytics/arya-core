package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryRetrieve", func() {
	BeforeEach(createDummyModel)
	AfterEach(deleteDummyModel)
	Describe("Retrieve an item", func() {
		It("Should retrieve it without error", func() {
			m := &storage.ChannelConfig{}
			err := dummyEngine.NewRetrieve(dummyAdapter).Model(m).WherePK(dummyModel.
				ID).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should retrieve the correct item", func() {
			m := &storage.ChannelConfig{}
			if err := dummyEngine.NewRetrieve(dummyAdapter).Model(m).WherePK(dummyModel.
				ID).Exec(dummyCtx); err != nil {
				log.Fatalln(err)
			}
			Expect(m.ID).To(Equal(dummyModel.ID))
			Expect(m.Name).To(Equal(dummyModel.Name))
		})
	})
	Describe("Retrieve multiple items", func() {
		It("Should retrieve all the correct items", func() {
			dummyModelTwo := &storage.ChannelConfig{
				ID:   uuid.New(),
				Name: "CC 45",
			}
			if err := dummyEngine.NewCreate(dummyAdapter).Model(dummyModelTwo).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}

			var models []*storage.ChannelConfig
			err := dummyEngine.NewRetrieve(dummyAdapter).Model(&models).WherePKs(
				[]uuid.UUID{dummyModelTwo.ID,
					dummyModel.ID}).Exec(dummyCtx)
			Expect(err).To(BeNil())
			Expect(models).To(HaveLen(2))
			Expect(models[0].Name == dummyModel.Name || models[0].
				Name == dummyModelTwo.Name).To(BeTrue())
			Expect(models[1].ID == dummyModelTwo.ID || models[1].ID == dummyModel.
				ID).To(BeTrue())
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct error type", func() {
				somePKThatDoesntExist := 136987
				m := &storage.ChannelConfig{}
				err := dummyEngine.NewRetrieve(dummyAdapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(dummyCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
			})
		})
	})
})
