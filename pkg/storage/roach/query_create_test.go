package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	var channelConfig *models.ChannelConfig
	var node *models.Node
	JustBeforeEach(func() {
		nErr := engine.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		ccErr := engine.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(ccErr).To(BeNil())
	})
	AfterEach(func() {
		ccErr := engine.NewDelete().Model(channelConfig).WherePK(channelConfig.
			ID).Exec(ctx)
		Expect(ccErr).To(BeNil())
		nErr := engine.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Create a new Channel Config", func() {
		BeforeEach(func() {
			node = &models.Node{ID: 1}
			channelConfig = &models.ChannelConfig{
				Name:   "Channel Config A",
				NodeID: node.ID,
			}
		})
		It("Should be able to be re-queried after creation", func() {
			resChannelConfig := &models.ChannelConfig{}
			err := engine.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.ID).
				Exec(ctx)
			Expect(err).To(BeNil())
			Expect(resChannelConfig).To(Equal(channelConfig))
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Unique Violation", func() {
			BeforeEach(func() {
				node = &models.Node{ID: 1}
				channelConfig = &models.ChannelConfig{
					Name:   "Channel Config A",
					NodeID: node.ID,
				}
			})
			It("Should return the correct errutil type", func() {
				channelConfigTwo := &models.ChannelConfig{
					ID:     channelConfig.ID,
					NodeID: node.ID,
				}
				err := engine.NewCreate().Model(channelConfigTwo).Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeUniqueViolation))
			})
		})
	})
})
