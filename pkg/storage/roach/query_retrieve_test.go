package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var channelConfig *models.ChannelConfig
	var node *models.Node
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New(), Name: "B"}
	})
	JustBeforeEach(func() {
		nErr := engine.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		ccErr := engine.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(ccErr).To(BeNil())
	})
	JustAfterEach(func() {
		ccErr := engine.NewDelete().Model(channelConfig).WherePK(channelConfig.
			ID).Exec(ctx)
		Expect(ccErr).To(BeNil())
		nErr := engine.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Describe("Retrieve an item", func() {
			It("Should retrieve it without error", func() {
				err := engine.NewRetrieve().Model(&models.ChannelConfig{}).
					WherePK(channelConfig.ID).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct item", func() {
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig).To(Equal(channelConfig))
			})
			It("Retrieve a single field", func() {
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve().Model(resChannelConfig).Fields("name").WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelConfig.ID).To(Equal(uuid.UUID{}))
				Expect(resChannelConfig.Name).To(Equal("B"))
			})
		})
		Describe("Retrieve multiple items", func() {
			var channelConfigTwo *models.ChannelConfig
			BeforeEach(func() {
				channelConfigTwo = &models.ChannelConfig{
					ID:     uuid.New(),
					Name:   "A",
					NodeID: 1,
				}
			})
			JustBeforeEach(func() {
				err := engine.NewCreate().Model(channelConfigTwo).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve all the correct items", func() {
				var models []*models.ChannelConfig
				err := engine.NewRetrieve().Model(&models).WherePKs(
					[]uuid.UUID{channelConfigTwo.ID, channelConfig.ID}).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect([]string{channelConfig.Name,
					channelConfigTwo.Name}).To(ContainElement(models[0].Name))
			})
			Describe("Limiting the number of results", func() {
				It("Should limit the number of results correctly", func() {
					var models []*models.ChannelConfig
					err := engine.NewRetrieve().
						Model(&models).
						WherePKs([]uuid.UUID{channelConfig.ID, channelConfigTwo.ID}).
						Limit(1).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(models).To(HaveLen(1))
				})
			})
			Describe("Ordering the results", func() {
				It("Should order the results correctly", func() {
					var models []*models.ChannelConfig
					err := engine.NewRetrieve().
						Model(&models).
						WherePKs([]uuid.UUID{channelConfig.ID, channelConfigTwo.ID}).
						Order(query.OrderASC, "Name").Exec(ctx)
					Expect(err).To(BeNil())
					Expect(models).To(HaveLen(2))
					Expect(models[0].Name).To(Equal("A"))
				})
				It("Should order the results by a nested field correctly", func() {
					var models []*models.ChannelConfig
					err := engine.NewRetrieve().
						Model(&models).
						WherePKs([]uuid.UUID{channelConfig.ID, channelConfigTwo.ID}).
						Relation("Node", "ID").
						Order(query.OrderASC, "Node.ID").Exec(ctx)
					Expect(err).To(BeNil())
					Expect(models).To(HaveLen(2))
					Expect(models[0].NodeID).To(Equal(1))
				})
			})
		})
		Describe("Retrieve a related item", func() {
			It("Should retrieve all of the correct items", func() {
				resChannelConfig := &models.ChannelConfig{}
				err := engine.NewRetrieve().Model(resChannelConfig).Relation("Node", "ID").
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
				rErr := engine.NewCreate().Model(rangeX).Exec(ctx)
				Expect(rErr).To(BeNil())
				ccErr := engine.NewCreate().Model(channelChunk).Exec(ctx)
				Expect(ccErr).To(BeNil())
				rrErr := engine.NewCreate().Model(rangeReplica).Exec(ctx)
				Expect(rrErr).To(BeNil())
				ccRErr := engine.NewCreate().Model(channelChunkReplica).Exec(ctx)
				Expect(ccRErr).To(BeNil())
			})
			It("Should retrieve all of the correct items", func() {
				channelChunkReplicaRes := &models.ChannelChunkReplica{}
				err := engine.NewRetrieve().Model(channelChunkReplicaRes).WherePK(channelChunkReplica.ID).Relation("RangeReplica.Node").Exec(ctx)
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
					err := engine.NewCreate().Model(item).Exec(ctx)
					Expect(err).To(BeNil())
				}
			})
			It("Should retrieve by the field correctly", func() {
				resRngLease := &models.RangeLease{}
				err := engine.
					NewRetrieve().
					Model(resRngLease).
					WhereFields(query.WhereFields{"RangeID": rng.ID}).
					Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resRngLease.ID).To(Equal(rngLease.ID))
			})
			It("Should retrieve the field by a great than expression", func() {

			})
			DescribeTable("Field Expressions", func(exp query.FieldExp) {
				cc := &models.ChannelChunk{ChannelConfigID: channelConfig.ID, Size: 4000, RangeID: rng.ID}
				cErr := engine.NewCreate().Model(cc).Exec(ctx)
				Expect(cErr).To(BeNil())
				var resCC []*models.ChannelChunk
				rErr := engine.NewRetrieve().Model(&resCC).WhereFields(query.WhereFields{"Size": exp}).Exec(ctx)
				Expect(rErr).To(BeNil())
				Expect(resCC).To(HaveLen(1))
			},
				Entry("Greater than", query.GreaterThan(3500)),
				Entry("Less than", query.LessThan(4500)),
				Entry("In Range", query.InRange(3000, 4500)),
			)

			It("Should return a not found error when no item can be found", func() {
				resRngLease := &models.RangeLease{}
				err := engine.
					NewRetrieve().
					Model(resRngLease).
					WhereFields(query.WhereFields{"RangeID": uuid.UUID{}}).
					Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
			Context("Nested Relation", func() {
				It("Should retrieve by a single nested relation correctly", func() {
					resRange := &models.Range{}
					err := engine.NewRetrieve().Model(resRange).WhereFields(query.WhereFields{
						"RangeLease.ID": rngLease.ID,
					}).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resRange.ID).To(Equal(rng.ID))

				})
				It("Should retrieve by a double nested relation correctly", func() {
					var resRanges []*models.Range
					err := engine.NewRetrieve().Model(&resRanges).WhereFields(query.WhereFields{
						"RangeLease.RangeReplica.NodeID": 1,
					}).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resRanges).To(HaveLen(1))
					Expect(resRanges[0].ID).To(Equal(rng.ID))
				})
			})

		})
		Describe("Using Calc", func() {
			var (
				size   int64 = 30
				count        = 10
				rng    *models.Range
				chunks []*models.ChannelChunk
			)
			BeforeEach(func() {
				rng = &models.Range{ID: uuid.New()}
				chunks = []*models.ChannelChunk{}
				for i := 0; i < count; i++ {
					chunks = append(chunks, &models.ChannelChunk{
						ID:              uuid.New(),
						RangeID:         rng.ID,
						ChannelConfigID: channelConfig.ID,
						Size:            size,
					})
				}
			})
			JustBeforeEach(func() {
				Expect(engine.NewCreate().Model(rng).Exec(ctx)).To(BeNil())
				Expect(engine.NewCreate().Model(&chunks).Exec(ctx)).To(BeNil())
			})
			JustAfterEach(func() {
				Expect(engine.NewDelete().Model(rng).WherePK(rng.ID).Exec(ctx)).To(BeNil())
				Expect(engine.NewDelete().Model(&chunks).WherePKs(model.NewReflect(&chunks).PKChain().Raw()).Exec(ctx)).To(BeNil())
			})
			Describe("Calculations", func() {
				It("Should calc the correct sum", func() {
					into := 0
					err := engine.NewRetrieve().
						Model(&models.ChannelChunk{}).
						Calc(query.CalcSum, "Size", &into).
						WherePKs(model.NewReflect(&chunks).PKChain().Raw()).
						Exec(ctx)
					Expect(err).To(BeNil())
					Expect(into).To(Equal(30 * 10))
				})
				It("Should calc the correct max", func() {
					into := 0
					err := engine.NewRetrieve().
						Model(&models.ChannelChunk{}).
						Calc(query.CalcMax, "Size", &into).
						WherePKs(model.NewReflect(&chunks).PKChain().Raw()).
						Exec(ctx)
					Expect(err).To(BeNil())
					Expect(into).To(Equal(30))
				})
				It("Should calc the correct min", func() {
					into := 0
					err := engine.NewRetrieve().
						Model(&models.ChannelChunk{}).
						Calc(query.CalcMin, "Size", &into).
						WherePKs(model.NewReflect(&chunks).PKChain().Raw()).
						Exec(ctx)
					Expect(err).To(BeNil())
					Expect(into).To(Equal(30))
				})
				It("Should calc the correct count", func() {
					into := 0
					err := engine.NewRetrieve().
						Model(&models.ChannelChunk{}).
						Calc(query.CalcCount, "Size", &into).
						WherePKs(model.NewReflect(&chunks).PKChain().Raw()).
						Exec(ctx)
					Expect(err).To(BeNil())
					Expect(into).To(Equal(10))
				})
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct errutil type", func() {
				somePKThatDoesntExist := uuid.New()
				m := &models.ChannelConfig{}
				err := engine.NewRetrieve().
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
	})
})
