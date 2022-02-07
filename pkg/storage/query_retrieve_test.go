package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var (
		node          *storage.Node
		channelConfig *storage.ChannelConfig
	)
	BeforeEach(func() {
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID,
			Name: "REALLY_AWESOME_SENSOR", ID: uuid.New()}
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
	Describe("Standard usage", func() {
		Context("Meta Data Only", func() {
			Context("Single item", func() {
				Describe("Retrieve a channel config", func() {
					It("Should retrieve the correct item", func() {
						resChannelConfig := &storage.ChannelConfig{}
						err := store.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.ID).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
						Expect(resChannelConfig.Name).To(Equal(channelConfig.Name))
					})
				})
			})
		})
		Context("Object Data + Meta Data", func() {
			Context("Single item", func() {
				var (
					channelChunk *storage.ChannelChunk
					bytes        []byte
				)
				BeforeEach(func() {
					bytes = []byte("randomstring")
					channelChunk = &storage.ChannelChunk{
						ChannelConfigID: channelConfig.ID,
						Data:            mock.NewObject(bytes),
					}
				})
				JustBeforeEach(func() {
					err := store.NewCreate().Model(channelChunk).Exec(ctx)
					Expect(err).To(BeNil())
				})
				JustAfterEach(func() {
					err := store.NewDelete().Model(channelChunk).WherePK(
						channelChunk.ID).Exec(ctx)
					Expect(err).To(BeNil())
				})
				Describe("Retrieve a channel chunk", func() {
					It("Should retrieve the correct item", func() {
						resChannelChunk := &storage.ChannelChunk{}
						err := store.NewRetrieve().Model(resChannelChunk).WherePK(
							channelChunk.ID).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelChunk.ID).To(Equal(channelChunk.ID))
						Expect(resChannelChunk.Data).ToNot(BeNil())
						resBytes := make([]byte, resChannelChunk.Data.Size())
						_, err = resChannelChunk.Data.Read(resBytes)
						Expect(err.Error()).To(Equal("EOF"))
						Expect(resBytes).To(Equal(bytes))
					})
				})
			})
		})
	})
})
