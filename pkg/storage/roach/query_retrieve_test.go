package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var channelConfig *models.ChannelConfig
	var node *models.Node
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New(), Name: "Channel Config"}
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
				err := engine.NewRetrieve(adapter).Model(&models.ChannelConfig{}).
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig).To(Equal(channelConfig))
			})
			It("Retrieve a single field", func() {
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).Fields("name").WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.ID).To(Equal(uuid.UUID{}))
				Expect(resChannelConfig.Name).To(Equal("Channel Config"))
			})
		})
		Describe("Retrieve multiple items", func() {
			var channelConfigTwo *models.ChannelConfig
			BeforeEach(func() {
				channelConfigTwo = &models.ChannelConfig{
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
				var models []*models.ChannelConfig
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
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve(adapter).Model(resChannelConfig).Relation("Node", "ID").
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.Node.ID).To(Equal(1))
			})
		})
		Describe("Retrieve through multiple levels of relations", func() {
			var (
				//rangeLease          *storage.RangeID
				rangeX              *models.Range
				channelChunkReplica *models.ChannelChunkReplica
				rangeReplica        *models.RangeReplica
				channelChunk        *models.ChannelChunk
			)
			BeforeEach(func() {
				rangeX = &models.Range{
					ID: uuid.New(),
				}
				channelChunk = &models.ChannelChunk{
					ID:              uuid.New(),
					RangeID:         rangeX.ID,
					ChannelConfigID: channelConfig.ID,
				}
				rangeReplica = &models.RangeReplica{
					ID:      uuid.New(),
					RangeID: rangeX.ID,
					NodeID:  node.ID,
				}
				channelChunkReplica = &models.ChannelChunkReplica{
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
				channelChunkReplicaRes := &models.ChannelChunkReplica{}
				err := engine.NewRetrieve(adapter).Model(channelChunkReplicaRes).WherePK(channelChunkReplica.ID).Relation("RangeReplica.Node").Exec(ctx)
				Expect(err).To(BeNil())
				Expect(channelChunkReplicaRes.RangeReplica.Node.ID).To(Equal(node.ID))

			})
		})
		Describe("Using WhereField", func() {
			var (
				rngLease   *models.RangeLease
				rng        *models.Range
				rngReplica *models.RangeReplica
				items      []interface{}
			)
			BeforeEach(func() {

				rng = &models.Range{
					ID: uuid.New(),
				}
				rngLease = &models.RangeLease{
					ID:      uuid.New(),
					RangeID: rng.ID,
				}
				rngReplica = &models.RangeReplica{
					ID:      uuid.New(),
					RangeID: rng.ID,
					NodeID:  node.ID,
				}
				rngLease.RangeReplicaID = rngReplica.ID
				items = []interface{}{
					rng,
					rngReplica,
					rngLease,
				}
			})
			JustBeforeEach(func() {
				for _, item := range items {
					err := engine.NewCreate(adapter).Model(item).Exec(ctx)
					Expect(err).To(BeNil())
				}
			})
			It("Should retrieve by the field correctly", func() {
				resRngLease := &models.RangeLease{}
				err := engine.
					NewRetrieve(adapter).
					Model(resRngLease).
					WhereFields(models.Fields{"RangeID": rng.ID}).
					Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resRngLease.ID).To(Equal(rngLease.ID))
			})
			It("Should return a not found error when no item can be found", func() {
				resRngLease := &models.RangeLease{}
				err := engine.
					NewRetrieve(adapter).
					Model(resRngLease).
					WhereFields(models.Fields{"RangeID": uuid.UUID{}}).
					Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
			Context("Nested Relation", func() {
				It("Should retrieve by a single nested relation correctly", func() {
					resRange := &models.Range{}
					err := engine.NewRetrieve(adapter).Model(resRange).WhereFields(models.Fields{
						"RangeLease.ID": rngLease.ID,
					}).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resRange.ID).To(Equal(rng.ID))

				})
				It("Should retrieve by a double nested relation correctly", func() {
					var resRanges []*models.Range
					err := engine.NewRetrieve(adapter).Model(&resRanges).WhereFields(models.Fields{
						"RangeLease.RangeReplica.NodeID": 1,
					}).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resRanges).To(HaveLen(1))
					Expect(resRanges[0].ID).To(Equal(rng.ID))
				})
				//It("Should return the correct error when an invalid relation is provided", func() {
				//	var resRanges []*models.Range
				//	err := engine.NewRetrieve(adapter).Model(&resRanges).WhereFields(models.Fields{
				//		"RangeLease.BadRel.NodeID": 1,
				//	}).Exec(ctx)
				//	Expect(err).ToNot(BeNil())
				//	Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeInvalidArgs))
				//})
			})

		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct errutil type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &models.ChannelConfig{}
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
