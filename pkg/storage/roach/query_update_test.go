package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryUpdate", func() {
	var channelConfig *storage.ChannelConfig
	var node *storage.Node
	BeforeEach(func() {
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
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
		var updatedChannelConfig *storage.ChannelConfig
		BeforeEach(func() {
			updatedChannelConfig = &storage.ChannelConfig{
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
			resChannelConfig := &storage.ChannelConfig{}
			err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.
				ID).Exec(ctx)
			Expect(err).To(BeNil())
			Expect(resChannelConfig.ID).To(Equal(updatedChannelConfig.ID))
			Expect(resChannelConfig.Name).To(Equal(updatedChannelConfig.Name))
		})
	})
})
