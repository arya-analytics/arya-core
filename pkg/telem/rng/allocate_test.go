package rng_test

//var _ = Describe("Rng", func() {
//	var (
//		node          *models.Node
//		rngItem       *models.Range
//		rngLease      *models.RangeLease
//		replica       *models.RangeReplica
//		channelConfig *models.ChannelConfig
//		chunks        []*models.ChannelChunk
//		items         []interface{}
//	)
//	BeforeEach(func() {
//		node = &models.Node{ID: 1}
//		channelConfig = &models.ChannelConfig{
//			ID:     uuid.New(),
//			NodeID: node.ID,
//		}
//		rngItem = &models.Range{
//			ID:   uuid.New(),
//			Open: true,
//		}
//
//		replica = &models.RangeReplica{
//			ID:      uuid.New(),
//			RangeID: rngItem.ID,
//			NodeID:  node.ID,
//		}
//		rngLease = &models.RangeLease{
//			RangeReplicaID: replica.ID,
//			RangeID:        rngItem.ID,
//		}
//		items = []interface{}{
//			node,
//			channelConfig,
//			rngItem,
//			replica,
//			rngLease,
//		}
//		for i := 0; i < 100; i++ {
//			chunks = append(chunks, &models.ChannelChunk{
//				ID:              uuid.New(),
//				ChannelConfigID: channelConfig.ID,
//				RangeID:         rngItem.ID,
//				Size:            models.ChannelChunkSizeLarge,
//			})
//		}
//		items = append(items, &chunks)
//	})
//	JustBeforeEach(func() {
//		for _, item := range items {
//			err := clust.NewCreate().Model(item).Exec(ctx)
//			if err != nil {
//				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeUniqueViolation))
//			}
//		}
//	})
//	Describe("Sizing", func() {
//		It("Should calculate the size of a range correctly", func() {
//			a := rng.Allocate{Cluster: clust}
//			rngSize := a.RetrieveRangeSize(ctx, model.NewReflect(rngItem).PK())
//			log.Infof("Range is %v percent full", (float64(rngSize)/float64(models.RangeSize))*100)
//		})
//	})
//	It("Allocating a chunk", func() {
//		a := rng.Allocate{Cluster: clust}
//		cc := &models.ChannelChunk{}
//		a.Chunk(ctx, 1, cc)
//		Expect(cc.RangeID).To(Equal(rngItem.ID))
//	})
//
//})
