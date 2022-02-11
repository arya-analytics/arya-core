package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var channelConfig *storage.ChannelConfig
	var node *storage.Node
	BeforeEach(func() {
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID}
	})
	JustBeforeEach(func() {
		nErr := engine.NewCreate(adapter).Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		ccErr := engine.NewCreate(adapter).Model(channelConfig).Exec(ctx)
		Expect(ccErr).To(BeNil())
	})
	AfterEach(func() {
		ccErr := engine.NewDelete(adapter).Model(channelConfig).WherePK(channelConfig.
			ID).Exec(ctx)
		Expect(ccErr).To(BeNil())
		nErr := engine.NewDelete(adapter).Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Describe("Retrieve an item", func() {
			It("Should retrieve it without error", func() {
				err := engine.NewRetrieve(adapter).Model(&storage.ChannelConfig{}).
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				resChannelConfig := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig).To(Equal(channelConfig))
			})
		})
		Describe("Retrieve multiple items", func() {
			var channelConfigTwo *storage.ChannelConfig
			BeforeEach(func() {
				channelConfigTwo = &storage.ChannelConfig{
					ID:     uuid.New(),
					Name:   "CC 45",
					NodeID: 1,
				}
			})
			JustBeforeEach(func() {
				err := engine.NewCreate(adapter).Model(channelConfigTwo).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve all the correct items", func() {
				var models []*storage.ChannelConfig
				err := engine.NewRetrieve(adapter).Model(&models).WherePKs(
					[]uuid.UUID{channelConfigTwo.ID,
						channelConfig.ID}).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect([]string{channelConfig.Name,
					channelConfigTwo.Name}).To(ContainElement(models[0].Name))
			})
		})
		Describe("Retrieve a related item", func() {
			It("Should retrieve all of the correct items", func() {
				resChannelConfig := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).Relation("Node").
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.Node.ID).To(Equal(1))

			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct errutil type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
	})
})
