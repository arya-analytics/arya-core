package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var channelConfig *storage.ChannelConfig
	var node *storage.Node
	BeforeEach(func() {
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New(), Name: "Channel Config"}
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
	Describe("Standard Usage", func() {
		Describe("Retrieve an item", func() {
			It("Should retrieve it without error", func() {
				err := engine.NewRetrieve(adapter).Model(&storage.ChannelConfig{}).
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				resChannelConfig := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig).To(Equal(channelConfig))
			})
			It("Retrieve a single field", func() {
				resChannelConfig := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).Field("name").WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.ID).To(Equal(uuid.UUID{}))
				Expect(resChannelConfig.Name).To(Equal("Channel Config"))
			})
		})
		Describe("Retrieve multiple items", func() {
			var channelConfigTwo *storage.ChannelConfig
			BeforeEach(func() {
				channelConfigTwo = &storage.ChannelConfig{
					ID:     uuid.New(),
					Name:   "CC 45",
					NodeID: 1,
				}
			})
			JustBeforeEach(func() {
				err := engine.NewCreate(adapter).Model(channelConfigTwo).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve all the correct items", func() {
				var models []*storage.ChannelConfig
				err := engine.NewRetrieve(adapter).Model(&models).WherePKs(
					[]uuid.UUID{channelConfigTwo.ID,
						channelConfig.ID}).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect([]string{channelConfig.Name,
					channelConfigTwo.Name}).To(ContainElement(models[0].Name))
			})
		})
		Describe("Retrieve a related item", func() {
			It("Should retrieve all of the correct items", func() {
				resChannelConfig := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).Relation("Node").
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.Node.ID).To(Equal(1))
			})
		})
		Describe("Retrieve through multiple levels of relations", func() {
			var (
				//rangeLease          *storage.RangeLease
				rangeX              *storage.Range
				channelChunkReplica *storage.ChannelChunkReplica
				rangeReplica        *storage.RangeReplica
				channelChunk        *storage.ChannelChunk
			)
			BeforeEach(func() {
				rangeX = &storage.Range{
					ID: uuid.New(),
				}
				channelChunk = &storage.ChannelChunk{
					ID:              uuid.New(),
					RangeID:         rangeX.ID,
					ChannelConfigID: channelConfig.ID,
				}
				rangeReplica = &storage.RangeReplica{
					ID:      uuid.New(),
					RangeID: rangeX.ID,
					NodeID:  node.ID,
				}
				channelChunkReplica = &storage.ChannelChunkReplica{
					RangeReplicaID: rangeReplica.ID,
					ChannelChunkID: channelChunk.ID,
				}
			})
			JustBeforeEach(func() {
				rErr := engine.NewCreate(adapter).Model(rangeX).Exec(ctx)
				Expect(rErr).To(BeNil())
				ccErr := engine.NewCreate(adapter).Model(channelChunk).Exec(ctx)
				Expect(ccErr).To(BeNil())
				rrErr := engine.NewCreate(adapter).Model(rangeReplica).Exec(ctx)
				Expect(rrErr).To(BeNil())
				ccRErr := engine.NewCreate(adapter).Model(channelChunkReplica).Exec(ctx)
				Expect(ccRErr).To(BeNil())
			})
			It("Should retrieve all of the correct items", func() {
				channelChunkReplicaRes := &storage.ChannelChunkReplica{}
				err := engine.NewRetrieve(adapter).Model(channelChunkReplicaRes).WherePK(channelChunkReplica.ID).Relation("RangeReplica.Node").Exec(ctx)
				Expect(err).To(BeNil())
				Expect(channelChunkReplicaRes.RangeReplica.Node.ID).To(Equal(node.ID))

			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct errutil type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &storage.ChannelConfig{}
				err := engine.NewRetrieve(adapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
	})
})
