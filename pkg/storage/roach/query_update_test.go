package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryUpdate", func() {
	var channelConfig *models.ChannelConfig
	var node *models.Node
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
	})
	JustBeforeEach(func() {
		nErr := engine.NewCreate(adapter).Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		ccErr := engine.NewCreate(adapter).Model(channelConfig).Exec(ctx)
		Expect(ccErr).To(BeNil())
	})
	JustAfterEach(func() {
		ccErr := engine.NewDelete(adapter).Model(channelConfig).WherePK(channelConfig.
			ID).Exec(ctx)
		Expect(ccErr).To(BeNil())
		nErr := engine.NewDelete(adapter).Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Update an item", func() {
		var updatedChannelConfig *models.ChannelConfig
		BeforeEach(func() {
			updatedChannelConfig = &models.ChannelConfig{
				ID:     channelConfig.ID,
				Name:   "Cool New Named Name",
				NodeID: 1,
			}
		})
		JustBeforeEach(func() {
			err := engine.NewUpdate(adapter).Model(updatedChannelConfig).WherePK(
				channelConfig.ID).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should reflect updates when retrieved", func() {
			resChannelConfig := &models.ChannelConfig{}
			err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.
				ID).Exec(ctx)
			Expect(err).To(BeNil())
			Expect(resChannelConfig.ID).To(Equal(updatedChannelConfig.ID))
			Expect(resChannelConfig.Name).To(Equal(updatedChannelConfig.Name))
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
			err := engine.NewCreate(adapter).Model(&channelConfigs).Exec(ctx)
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
			err = engine.NewUpdate(adapter).Model(&updateConfigs).Fields("Name").Bulk().Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should reflect the updates when retrieved", func() {
			var resChannelConfigs []*models.ChannelConfig
			err := engine.NewRetrieve(adapter).
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
