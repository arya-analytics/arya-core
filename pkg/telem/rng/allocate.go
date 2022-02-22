package rng

//type Allocate struct {
//	TrackOpen Track
//	Cluster   cluster.Cluster
//}
//
//func (a *Allocate) Chunk(ctx context.Context, nodeID int, chunk *models.ChannelChunk) {
//	rangeID, _, ok := a.TrackOpen.RetrieveOpen(nodeID)
//	if ok {
//		chunk.RangeID = rangeID
//		return
//	}
//	r := a.Range(ctx, nodeID)
//	chunk.RangeID = r.ID
//}
//
//func (a *Allocate) ChunkReplica(ctx context.Context, rangeID uuid.UUID, replica *models.ChannelChunkReplica) {
//	rangeInfo, ok := a.TrackOpen.Retrieve(rangeID)
//	if !ok || !rangeInfo.Open {
//		r := &models.Range{}
//		if err := a.Cluster.NewRetrieve().
//			Model(r).
//			WherePK(rangeID).
//			Relation("RangeLease.RangeReplica", "ID, NodeID").
//			Exec(ctx); err != nil {
//			panic(err)
//		}
//		rangeInfo.Open = r.Open
//		rangeInfo.LeaseNodeID = r.RangeLease.RangeReplica.NodeID
//		rangeInfo.LeaseReplicaID = r.RangeLease.RangeReplica.ID
//		a.TrackOpen.Add(r.ID, rangeInfo.Open, rangeInfo.LeaseNodeID, rangeInfo.LeaseReplicaID)
//	}
//	if !rangeInfo.Open {
//		r := a.Range(ctx, rangeInfo.LeaseNodeID)
//		rangeInfo.LeaseReplicaID = r.RangeLease.RangeReplicaID
//	}
//	replica.RangeReplicaID = rangeInfo.LeaseReplicaID
//}
//
//func (a *Allocate) Range(ctx context.Context, nodeID int) *models.Range {
//	r := &models.Range{}
//	if err := a.Cluster.NewCreate().Model(r).Exec(ctx); err != nil {
//		panic(err)
//	}
//	rr := &models.RangeReplica{
//		RangeID: r.ID,
//		NodeID:  nodeID,
//	}
//	if err := a.Cluster.NewCreate().Model(rr).Exec(ctx); err != nil {
//		panic(err)
//	}
//	lease := &models.RangeLease{
//		RangeID:        r.ID,
//		RangeReplicaID: rr.ID,
//	}
//	if err := a.Cluster.NewCreate().Model(lease).Exec(ctx); err != nil {
//		panic(err)
//	}
//	lease.RangeReplica = rr
//	r.RangeLease = lease
//	a.TrackOpen.Add(r.ID, r.Open, rr.NodeID, rr.ID)
//	return r
//}
//
//func (a *Allocate) RetrieveRangeSize(ctx context.Context, pk model.PK) (rangeSize int64) {
//	var chunks []*models.ChannelChunk
//	if err := a.Cluster.NewRetrieve().Model(&chunks).WhereFields(model.WhereFields{"RangeID": pk.Raw()}).Exec(ctx); err != nil {
//		panic(err)
//	}
//	for _, chunk := range chunks {
//		rangeSize += int64(chunk.Size)
//	}
//	return rangeSize
//
//}
