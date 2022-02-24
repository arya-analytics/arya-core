package chanchunk

//
//type Service struct {
//	cluster cluster.Cluster
//	alloc   *rng.Allocate
//}
//
//func (s *Service) CreateStream(ctx context.Context, cfg *models.ChannelConfig) error {
//	channelConfig := &models.ChannelConfig{}
//	if err := s.cluster.NewRetrieve().
//		Model(channelConfig).
//		WherePK(stream.ChannelConfigID).
//		Exec(ctx); err != nil {
//		return err
//	}
//	for {
//		select {
//		case repl, ok := <-stream.ChunkReplicas:
//			if !ok {
//				return nil
//			}
//			chunk := &models.ChannelChunk{PK: uuid.New(), ChannelConfigID: channelConfig.PK}
//			s.alloc.Chunk(ctx, channelConfig.NodeID, chunk)
//			if err := s.cluster.NewCreate().Model(chunk).Exec(ctx); err != nil {
//				return err
//			}
//			s.alloc.ChunkReplica(ctx, chunk.RangeID, repl)
//			if err := s.cluster.NewCreate().Model(repl).Exec(ctx); err != nil {
//				return err
//			}
//		}
//	}
//}
