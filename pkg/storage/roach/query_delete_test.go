package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryDelete", func() {
	var channelConfig *models.ChannelConfig
	var node *models.Node
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID}
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
	Describe("Delete an item", func() {
		It("Should delete correctly", func() {
			dErr := engine.NewDelete(adapter).Model(channelConfig).WherePK(
				channelConfig.ID).Exec(ctx)
			Expect(dErr).To(BeNil())
			rErr := engine.NewRetrieve(adapter).Model(channelConfig).WherePK(channelConfig.ID).
				Exec(ctx)
			Expect(rErr).ToNot(BeNil())
			Expect(rErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
		})
	})
	Describe("Delete multiple items", func() {
		var channelConfigTwo *models.ChannelConfig
		BeforeEach(func() {
			channelConfigTwo = &models.ChannelConfig{
				ID:     uuid.New(),
				Name:   "CC 45",
				NodeID: node.ID,
			}
		})
		JustBeforeEach(func() {
			cErr := engine.NewCreate(adapter).Model(channelConfigTwo).Exec(ctx)
			Expect(cErr).To(BeNil())
		})
		It("Should delete them correctly", func() {
			pks := []uuid.UUID{channelConfig.ID, channelConfigTwo.ID}
			err := engine.NewDelete(adapter).Model(&models.ChannelConfig{}).
				WherePKs(pks).
				Exec(ctx)
			Expect(err).To(BeNil())
			var models []*models.ChannelConfig
			e := engine.NewRetrieve(adapter).Model(&models).WherePKs(
				[]uuid.UUID{channelConfigTwo.ID,
					channelConfig.ID}).Exec(ctx)
			Expect(e).To(BeNil())
			Expect(models).To(HaveLen(0))
		})
	})
})
