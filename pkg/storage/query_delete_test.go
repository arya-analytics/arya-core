package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryDelete", func() {
	var (
		node          *models.Node
		channelConfig *models.ChannelConfig
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID}
	})
	JustBeforeEach(func() {
		nErr := store.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		cErr := store.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(cErr).To(BeNil())
	})
	JustAfterEach(func() {
		nErr := store.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Describe("Delete a channel config", func() {
			It("Should delete correctly", func() {
				dErr := store.NewDelete().Model(channelConfig).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(dErr).To(BeNil())
				rErr := store.NewRetrieve().Model(channelConfig).WherePK(
					channelConfig.ID).Exec(ctx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Describe("Providing no where statement to the query", func() {
			It("Should return an error", func() {
				err := store.NewDelete().Model(&models.ChannelConfig{}).Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeInvalidArgs))
			})
		})
	})
})
