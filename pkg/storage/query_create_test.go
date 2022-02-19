package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	var (
		node           *models.Node
		channelConfig  *models.ChannelConfig
		channelConfigs []*models.ChannelConfig
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
	})
	JustBeforeEach(func() {
		err := store.NewCreate().Model(node).Exec(ctx)
		Expect(err).To(BeNil())
	})
	JustAfterEach(func() {
		err := store.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(err).To(BeNil())
	})
	Describe("Create a  new item", func() {
		BeforeEach(func() {
			channelConfig = &models.ChannelConfig{
				NodeID: node.ID,
			}
		})
		JustAfterEach(func() {
			err := store.NewDelete().Model(channelConfig).WherePK(channelConfig.ID).
				Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should create the correct item", func() {
			err := store.NewCreate().Model(channelConfig).Exec(ctx)
			Expect(err).To(BeNil())
			resChannelConfig := &models.ChannelConfig{}
			rErr := store.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.
				ID).Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resChannelConfig).To(Equal(channelConfig))
		})
	})
	Describe("Create multiple items", func() {
		BeforeEach(func() {
			channelConfigs = []*models.ChannelConfig{
				{
					Name:   "Cool Name 1",
					NodeID: node.ID,
				},
				{
					Name:   "Cool Name 2",
					NodeID: node.ID,
				},
			}
		})
		JustAfterEach(func() {
			pks := []uuid.UUID{channelConfigs[0].ID, channelConfigs[1].ID}
			err := store.NewDelete().Model(&channelConfigs).WherePKs(pks).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should create the items correctly", func() {
			var resChannelConfigs []*models.ChannelConfig
			By("Creating without error")
			cErr := store.NewCreate().Model(&channelConfigs).Exec(ctx)
			Expect(cErr).To(BeNil())
			By("Being re-retrieved after creation")
			pks := []uuid.UUID{channelConfigs[0].ID, channelConfigs[1].ID}
			rErr := store.NewRetrieve().Model(&resChannelConfigs).WherePKs(pks).
				Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resChannelConfigs).To(HaveLen(2))
			Expect(pks).To(ContainElements(resChannelConfigs[0].ID,
				resChannelConfigs[1].ID))
		})
	})
})
