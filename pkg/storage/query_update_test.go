package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
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
	Describe("Bulk Update Items", func() {
		var (
			channelConfigs []*models.ChannelConfig
		)
		BeforeEach(func() {
			channelConfigs = []*models.ChannelConfig{
				{
					Name:     "Hello",
					NodeID:   node.ID,
					DataRate: 32,
				},
				{
					Name:     "Hello 2",
					NodeID:   node.ID,
					DataRate: 32,
				},
			}
		})
		JustBeforeEach(func() {
			err := store.NewCreate().Model(&channelConfigs).Exec(ctx)
			Expect(err).To(BeNil())
			updateConfigs := []*models.ChannelConfig{
				{
					ID:   channelConfigs[0].ID,
					Name: "New Name",
				},
				{
					ID:   channelConfigs[1].ID,
					Name: "New Name",
				},
			}
			err = store.NewUpdate().Model(&updateConfigs).Fields("Name").Bulk().Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should reflect the updates when retrieved", func() {
			var resChannelConfigs []*models.ChannelConfig
			err := store.NewRetrieve().
				Model(&resChannelConfigs).
				WherePKs(model.NewReflect(&channelConfigs).PKChain().Raw()).
				Exec(ctx)
			Expect(err).To(BeNil())
			Expect(len(resChannelConfigs)).To(Equal(2))
			Expect(resChannelConfigs[0].ID).To(BeElementOf(channelConfigs[0].ID, channelConfigs[1].ID))
			Expect(resChannelConfigs[0].Name).To(Equal("New Name"))
			Expect(resChannelConfigs[1].Name).To(Equal("New Name"))
			Expect(resChannelConfigs[0].DataRate).To(Equal(32))
			Expect(resChannelConfigs[1].DataRate).To(Equal(32))
		})
	})
})
