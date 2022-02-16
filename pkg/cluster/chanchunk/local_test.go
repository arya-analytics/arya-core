package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/storage"
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
		channelConfig       *storage.ChannelConfig
		node                *storage.Node
		rangeX              *storage.Range
		channelChunkReplica *storage.ChannelChunkReplica
		rangeReplica        *storage.RangeReplica
		channelChunk        *storage.ChannelChunk
	)
	BeforeEach(func() {
		localSvc = chanchunk.NewServiceLocalStorage(store)
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rangeX = &storage.Range{
			ID: uuid.New(),
		}
		rangeReplica = &storage.RangeReplica{
			ID:      uuid.New(),
			RangeID: rangeX.ID,
			NodeID:  node.ID,
		}
		channelChunk = &storage.ChannelChunk{
			ID:              uuid.New(),
			RangeID:         rangeX.ID,
			ChannelConfigID: channelConfig.ID,
		}
		channelChunkReplica = &storage.ChannelChunkReplica{
			RangeReplicaID: rangeReplica.ID,
			ChannelChunkID: channelChunk.ID,
			Telem:          telem.NewBulk([]byte{}),
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
				ccErr := localSvc.Create(ctx, model.NewReflect(channelChunk))
				Expect(ccErr).To(BeNil())
			})
			It("Should create the correct chunk", func() {
				resCC := &storage.ChannelChunk{}
				err := localSvc.Retrieve(ctx, model.NewReflect(resCC), model.NewPKChain([]uuid.UUID{channelChunk.ID}))
				Expect(err).To(BeNil())
				Expect(resCC.ID).To(Equal(channelChunk.ID))
			})
		})
		Context("Channel Chunk Replica", func() {
			JustBeforeEach(func() {
				ccErr := localSvc.Create(ctx, model.NewReflect(channelChunk))
				Expect(ccErr).To(BeNil())
				ccrErr := localSvc.CreateReplicas(ctx, model.NewReflect(channelChunkReplica))
				Expect(ccrErr).To(BeNil())
			})
			JustAfterEach(func() {
				ccErr := localSvc.Delete(ctx, model.NewReflect(channelChunk).PKChain())
				Expect(ccErr).To(BeNil())
				ccrErr := localSvc.DeleteReplicas(ctx, model.NewReflect(channelChunkReplica).PKChain())
				Expect(ccrErr).To(BeNil())
			})
			It("Should create the  replica correctly", func() {
				resCCR := &storage.ChannelChunkReplica{}
				err := localSvc.RetrieveReplicas(ctx, model.NewReflect(resCCR), model.NewPKChain([]uuid.UUID{channelChunkReplica.ID}), false)
				Expect(err).To(BeNil())
				Expect(resCCR.ID).To(Equal(channelChunkReplica.ID))
			})
		})
		Context("Retrieve Range Replicas", func() {
			It("Should retrieve the node information correctly", func() {
				resRR := &storage.RangeReplica{}
				err := localSvc.RetrieveRangeReplicas(ctx, model.NewReflect(resRR), model.NewPKChain([]uuid.UUID{rangeReplica.ID}))
				Expect(err).To(BeNil())
				Expect(resRR.Node.ID).To(Equal(1))
			})
		})
	})
})
