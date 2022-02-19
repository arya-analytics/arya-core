package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update QueryRequest", func() {
	var (
		channelConfig *models.ChannelConfig
		node          *models.Node
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
	})
	JustBeforeEach(func() {
		nErr := store.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		cErr := store.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(cErr).To(BeNil())
	})
	JustAfterEach(func() {
		cErr := store.NewDelete().Model(channelConfig).WherePK(channelConfig.ID).
			Exec(ctx)
		Expect(cErr).To(BeNil())
		nErr := store.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Update an item", func() {
		var updateChannelConfig *models.ChannelConfig
		BeforeEach(func() {
			updateChannelConfig = &models.ChannelConfig{
				ID:     channelConfig.ID,
				Name:   "Cool Name",
				NodeID: node.ID,
			}
		})
		It("Should update it correctly", func() {
			uErr := store.NewUpdate().Model(updateChannelConfig).WherePK(channelConfig.ID).Exec(ctx)
			Expect(uErr).To(BeNil())
			resChannelConfig := &models.ChannelConfig{}
			rErr := store.NewRetrieve().Model(resChannelConfig).WherePK(
				updateChannelConfig.ID).Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
			Expect(resChannelConfig.ID).To(Equal(updateChannelConfig.ID))
			Expect(resChannelConfig.Name).To(Equal(updateChannelConfig.Name))
		})
	})
	Describe("Edge cases + errors", func() {
		It("Should panic if a chain of models is provided", func() {
			Expect(func() { store.NewUpdate().Model(&[]*models.ChannelConfig{}) }).To(Panic())
		})
	})
})
