package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	var channelConfig *storage.ChannelConfig
	var node *storage.Node
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
	Describe("Create a new Channel Config", func() {
		BeforeEach(func() {
			node = &storage.Node{ID: 1}
			channelConfig = &storage.ChannelConfig{
				Name:   "Channel Config",
				NodeID: node.ID,
			}
		})
		It("Should be able to be re-queried after creation", func() {
			resChannelConfig := &storage.ChannelConfig{}
			err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.ID).
				Exec(ctx)
			Expect(err).To(BeNil())
			Expect(resChannelConfig).To(Equal(channelConfig))
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Unique Violation", func() {
			BeforeEach(func() {
				node = &storage.Node{ID: 1}
				channelConfig = &storage.ChannelConfig{
					Name:   "ChannelConfig",
					NodeID: node.ID,
				}
			})
			It("Should return the correct errutil type", func() {
				channelConfigTwo := &storage.ChannelConfig{
					ID:     channelConfig.ID,
					NodeID: node.ID,
				}
				err := engine.NewCreate(adapter).Model(channelConfigTwo).Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeUniqueViolation))
			})
		})
	})
})
