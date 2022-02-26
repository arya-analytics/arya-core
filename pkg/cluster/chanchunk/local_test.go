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
		localSvc            chanchunk.Local
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
			channelChunk,
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
	Describe("Channel Chunk Replica", func() {
		JustBeforeEach(func() {
			ccrErr := localSvc.Create(ctx, channelChunkReplica)
			Expect(ccrErr).To(BeNil())
		})
		JustAfterEach(func() {
			ccrOpts := chanchunk.LocalDeleteOpts{
				PKC: model.NewPKChain([]uuid.UUID{channelChunk.ID}),
			}
			ccrErr := localSvc.Delete(ctx, ccrOpts)
			Expect(ccrErr).To(BeNil())
		})
		It("Should create the replica correctly", func() {
			resCCR := &models.ChannelChunkReplica{}
			opts := chanchunk.LocalRetrieveOpts{
				PKC: model.NewPKChain([]uuid.UUID{channelChunkReplica.ID}),
			}
			err := localSvc.Retrieve(ctx, resCCR, opts)
			Expect(err).To(BeNil())
			Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
			Expect(resCCR.Telem.Bytes()).To(Equal([]byte("randomdata")))
		})
		It("Should omit bulk data when option specified", func() {
			resCCR := &models.ChannelChunkReplica{}
			opts := chanchunk.LocalRetrieveOpts{
				PKC:       model.NewPKChain([]uuid.UUID{channelChunkReplica.ID}),
				OmitTelem: true,
			}
			err := localSvc.Retrieve(ctx, resCCR, opts)
			Expect(err).To(BeNil())
			Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
			Expect(resCCR.Telem).To(BeNil())
		})
		It("Should bulk update the replica correctly", func() {
			newRR := &models.RangeReplica{ID: uuid.New(), RangeID: rangeX.ID, NodeID: node.ID}
			Expect(store.NewCreate().Model(newRR).Exec(ctx)).To(BeNil())
			updateCCR := []*models.ChannelChunkReplica{{ID: channelChunkReplica.ID, RangeReplicaID: newRR.ID}}
			opts := chanchunk.LocalUpdateOpts{
				Bulk:   true,
				Fields: []string{"RangeReplicaID"},
			}
			Expect(localSvc.Update(ctx, &updateCCR, opts)).To(BeNil())
		})
	})
	Describe("Retrieve Range Replicas", func() {
		It("Should retrieve the node information correctly", func() {
			resRR := &models.RangeReplica{}
			err := localSvc.RetrieveRangeReplica(ctx, resRR, model.NewPKChain([]uuid.UUID{rangeReplica.ID}))
			Expect(err).To(BeNil())
			Expect(resRR.Node.ID).To(Equal(1))
		})
	})
})
