package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local", func() {
	var (
		localSvc            chanchunk.ServiceLocal
		items               []interface{}
		channelConfig       *models.ChannelConfig
		node                *models.Node
		rangeX              *models.Range
		channelChunkReplica *models.ChannelChunkReplica
		rangeReplica        *models.RangeReplica
		channelChunk        *models.ChannelChunk
	)
	BeforeEach(func() {
		localSvc = chanchunk.NewServiceLocalStorage(store)
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rangeX = &models.Range{
			ID: uuid.New(),
		}
		rangeReplica = &models.RangeReplica{
			ID:      uuid.New(),
			RangeID: rangeX.ID,
			NodeID:  node.ID,
		}
		channelChunk = &models.ChannelChunk{
			ID:              uuid.New(),
			RangeID:         rangeX.ID,
			ChannelConfigID: channelConfig.ID,
		}
		channelChunkReplica = &models.ChannelChunkReplica{
			RangeReplicaID: rangeReplica.ID,
			ChannelChunkID: channelChunk.ID,
			Telem:          telem.NewBulk([]byte("randomdata")),
		}
		items = []interface{}{
			node,
			channelConfig,
			rangeX,
			rangeReplica,
		}
	})
	JustBeforeEach(func() {
		for _, m := range items {
			err := store.NewCreate().Model(m).Exec(ctx)
			Expect(err).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, m := range items {
			err := store.NewDelete().Model(m).WherePK(model.NewReflect(m).PK().Raw()).Exec(ctx)
			Expect(err).To(BeNil())
		}
	})
	Describe("Local storage", func() {
		Context("Channel Chunk", func() {
			JustBeforeEach(func() {
				ccErr := localSvc.CreateChunk(ctx, channelChunk)
				Expect(ccErr).To(BeNil())
			})
			It("Should create the correct chunk", func() {
				resCC := &models.ChannelChunk{}
				opts := chanchunk.LocalChunkRetrieveOpts{
					PKC: model.NewPKChain([]uuid.UUID{channelChunk.ID}),
				}
				err := localSvc.RetrieveChunk(ctx, resCC, opts)
				Expect(err).To(BeNil())
				Expect(resCC.ID).To(Equal(channelChunk.ID))
			})
		})
		Context("Channel Chunk Replica", func() {
			JustBeforeEach(func() {
				ccErr := localSvc.CreateChunk(ctx, channelChunk)
				Expect(ccErr).To(BeNil())
				ccrErr := localSvc.CreateReplica(ctx, channelChunkReplica)
				Expect(ccrErr).To(BeNil())
			})
			JustAfterEach(func() {
				ccOpts := chanchunk.LocalChunkDeleteOpts{
					PKC: model.NewPKChain([]uuid.UUID{channelChunk.ID}),
				}
				ccErr := localSvc.DeleteChunk(ctx, ccOpts)
				Expect(ccErr).To(BeNil())
				ccrOpts := chanchunk.LocalReplicaDeleteOpts{
					PKC: model.NewPKChain([]uuid.UUID{channelChunk.ID}),
				}
				ccrErr := localSvc.DeleteReplica(ctx, ccrOpts)
				Expect(ccrErr).To(BeNil())
			})
			It("Should create the replica correctly", func() {
				resCCR := &models.ChannelChunkReplica{}
				opts := chanchunk.LocalReplicaRetrieveOpts{
					PKC: model.NewPKChain([]uuid.UUID{channelChunkReplica.ID}),
				}
				err := localSvc.RetrieveReplica(ctx, resCCR, opts)
				Expect(err).To(BeNil())
				Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
				Expect(resCCR.Telem.Bytes()).To(Equal([]byte("randomdata")))
			})
			It("Should omit bulk data when option specified", func() {
				resCCR := &models.ChannelChunkReplica{}
				opts := chanchunk.LocalReplicaRetrieveOpts{
					PKC:      model.NewPKChain([]uuid.UUID{channelChunkReplica.ID}),
					OmitBulk: true,
				}
				err := localSvc.RetrieveReplica(ctx, resCCR, opts)
				Expect(err).To(BeNil())
				Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
				Expect(resCCR.Telem).To(BeNil())
			})
		})
		Context("Retrieve Range Replicas", func() {
			It("Should retrieve the node information correctly", func() {
				resRR := &models.RangeReplica{}
				opts := chanchunk.LocalRangeReplicaRetrieveOpts{
					PKC: model.NewPKChain([]uuid.UUID{rangeReplica.ID}),
				}
				err := localSvc.RetrieveRangeReplica(ctx, resRR, opts)
				Expect(err).To(BeNil())
				Expect(resRR.Node.ID).To(Equal(1))
			})
		})
	})
})
