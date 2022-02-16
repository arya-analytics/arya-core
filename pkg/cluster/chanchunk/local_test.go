package chanchunk_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Local", func() {
	//var (
	//	svc                 *chanchunk.ServiceLocal
	//	local               localstorage.Local
	//	items               []interface{}
	//	channelConfig       *storage.ChannelConfig
	//	node                *storage.Node
	//	rangeX              *storage.Range
	//	channelChunkReplica *storage.ChannelChunkReplica
	//	rangeReplica        *storage.RangeReplica
	//	channelChunk        *storage.ChannelChunk
	//)
	//BeforeEach(func() {
	//	local = &local.LocalStorage{
	//		Storage: chanchunk.store,
	//	}
	//	svc = &chanchunk.ServiceLocal{
	//		Local:   local,
	//		Catcher: &errutil.Catcher{},
	//	}
	//	node = &storage.Node{ID: 1}
	//	channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
	//	rangeX = &storage.Range{
	//		ID: uuid.New(),
	//	}
	//	channelChunk = &storage.ChannelChunk{
	//		ID:              uuid.New(),
	//		RangeID:         rangeX.ID,
	//		ChannelConfigID: channelConfig.ID,
	//	}
	//	rangeReplica = &storage.RangeReplica{
	//		ID:      uuid.New(),
	//		RangeID: rangeX.ID,
	//		NodeID:  node.ID,
	//	}
	//	channelChunkReplica = &storage.ChannelChunkReplica{
	//		RangeReplicaID: rangeReplica.ID,
	//		ChannelChunkID: channelChunk.ID,
	//		Telem:          telem.NewBulk([]byte{}),
	//	}
	//	items = []interface{}{
	//		node,
	//		channelConfig,
	//		rangeX,
	//		channelChunk,
	//		rangeReplica,
	//		channelChunkReplica,
	//	}
	//})
	//JustBeforeEach(func() {
	//	for _, m := range items {
	//		err := chanchunk.store.NewCreate().Model(m).Model(m).Send(chanchunk.ctx)
	//		Expect(err)
	//	}
	//})
	//Describe("Local storage", func() {
	//	It("Creating a channel chunk", func() {
	//		blk := telem.NewBulk([]byte{})
	//		mock.TelemBulkPopulateRandomFloat64(blk, 1000000)
	//		createCCR := &storage.ChannelChunkReplica{
	//			RangeReplicaID: rangeReplica.ID,
	//			ChannelChunkID: channelChunk.ID,
	//			Telem:          blk,
	//		}
	//		q := &cluster.QueryRequest{
	//			Model:   model.NewReflect(createCCR),
	//			Variant: cluster.QueryVariantCreate,
	//		}
	//		svc.Create(chanchunk.ctx, q)
	//	})
	//	Context("Retrieving channel chunk nodes", func() {
	//		//It("Should retrieve the correct node", func() {
	//		//	ccRfl := model.NewReflect(channelChunkReplica)
	//		//	nodes, err := Local.RetrieveRangeReplicas(ctx, ccRfl, ccRfl.PKChain())
	//		//	Expect(err).To(BeNil())
	//		//	Expect(nodes).To(HaveLen(1))
	//		//	Expect(nodes[0].ID).To(Equal(1))
	//		//	Expect(len(nodes[0].Address)).To(BeNumerically(">", 9))
	//		//	Expect(nodes[0].IsHost).To(BeTrue())
	//		//})
	//	})
	//})

})
