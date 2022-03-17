package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hook", func() {
	var node *models.Node
	BeforeEach(func() {
		node = &models.Node{ID: 1}
	})
	JustBeforeEach(func() {
		nErr := engine.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	AfterEach(func() {
		nErr := engine.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("UUID auto-generation", func() {
		var channelConfig *models.ChannelConfig
		BeforeEach(func() {
			channelConfig = &models.ChannelConfig{Name: "Auto-generated UUID", NodeID: node.ID}
		})
		JustBeforeEach(func() {
			Expect(engine.NewCreate().Model(channelConfig).Exec(ctx)).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			var retrievedCC = &models.ChannelConfig{}
			Expect(engine.NewRetrieve().Model(retrievedCC).WherePK(channelConfig.ID).Exec(ctx)).To(BeNil())
			Expect(retrievedCC.Name).To(Equal(channelConfig.Name))
		})
	})
})
