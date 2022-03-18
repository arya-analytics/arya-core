package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var (
		node          *models.Node
		channelConfig *models.ChannelConfig
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID,
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
		Context("Meta Telem Only", func() {
			Context("Single item", func() {
				Describe("Acquire a channel config", func() {
					It("Should Acquire the correct item", func() {
						resChannelConfig := &models.ChannelConfig{}
						err := store.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.ID).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
						Expect(resChannelConfig.Name).To(Equal(channelConfig.Name))
					})
					It("Should Acquire the channel config by a relation", func() {
						resChannelConfig := &models.ChannelConfig{}
						err := store.NewRetrieve().Model(resChannelConfig).WhereFields(query.WhereFields{
							"Node.ID": 1,
						}).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
						Expect(resChannelConfig.Name).To(Equal(channelConfig.Name))
					})
					It("Should Acquire the correct relation", func() {
						resChannelConfig := &models.ChannelConfig{}
						err := store.NewRetrieve().Model(resChannelConfig).Relation("Node", "id").WhereFields(query.WhereFields{
							"Node.ID": 1,
						}).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
						Expect(resChannelConfig.Name).To(Equal(channelConfig.Name))
						Expect(resChannelConfig.Node.ID).To(Equal(1))
					})
					It("Should Acquire only the specified fields", func() {
						resChannelConfig := &models.ChannelConfig{}
						err := store.NewRetrieve().Model(resChannelConfig).WherePK(channelConfig.ID).Fields("ID").Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resChannelConfig.ID).To(Equal(channelConfig.ID))
						Expect(resChannelConfig.Name).ToNot(Equal(channelConfig.Name))
					})
					It("Should run a max calculation", func() {
						max := 0
						resChannelConfig := &models.ChannelConfig{}
						err := store.NewRetrieve().
							Model(resChannelConfig).
							WherePK(channelConfig.ID).
							Calc(query.CalcMax, "NodeID", &max).
							Exec(ctx)
						Expect(err).To(BeNil())
						Expect(max).To(Equal(1))
					})
				})
			})
		})
		Context("Object + Meta Data", func() {
			Context("Single item", func() {
				var (
					channelChunk        *models.ChannelChunk
					channelChunkReplica *models.ChannelChunkReplica
					bytes               []byte
					items               []interface{}
				)
				BeforeEach(func() {
					bytes = []byte("randomstring")
					rng := &models.Range{ID: uuid.New()}
					rngReplica := &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
					channelChunk = &models.ChannelChunk{
						ID:              uuid.New(),
						RangeID:         rng.ID,
						ChannelConfigID: channelConfig.ID,
					}
					channelChunkReplica = &models.ChannelChunkReplica{
						ChannelChunkID: channelChunk.ID,
						Telem:          telem.NewChunkData(bytes),
						RangeReplicaID: rngReplica.ID,
					}
					items = []interface{}{
						rng,
						rngReplica,
						channelChunk,
						channelChunkReplica,
					}
				})
				JustBeforeEach(func() {
					for _, item := range items {
						err := store.NewCreate().Model(item).Exec(ctx)
						Expect(err).To(BeNil())
					}
				})
				JustAfterEach(func() {
					for _, item := range items {
						err := store.NewDelete().Model(item).WherePK(model.NewReflect(item).PK().Raw()).Exec(ctx)
						Expect(err).To(BeNil())
					}
				})
				Describe("Acquire a channel chunk", func() {
					It("Should Acquire the correct item", func() {
						resCCR := &models.ChannelChunkReplica{}
						err := store.NewRetrieve().Model(resCCR).WherePK(
							channelChunkReplica.ID).Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
						Expect(resCCR.Telem).ToNot(BeNil())
						Expect(resCCR.Telem.Bytes()).To(Equal(bytes))
					})
				})
				Describe("Retrieving the chunk by a non pk field", func() {
					It("Should Acquire the correct item", func() {
						resCCR := &models.ChannelChunkReplica{}
						err := store.NewRetrieve().
							Model(resCCR).
							WhereFields(query.WhereFields{"ChannelChunkID": channelChunk.ID}).
							Exec(ctx)
						Expect(err).To(BeNil())
						Expect(resCCR.Telem.Bytes()).To(Equal(bytes))
					})
				})
			})
		})
	})
})
