package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

type Persist interface {
	NewRange(ctx context.Context, nodeID int) (*models.Range, error)
	NewRangeReplica(ctx context.Context, rangeID uuid.UUID, nodeID int) (*models.RangeReplica, error)

	RetrieveRange(ctx context.Context, ID uuid.UUID) (*models.Range, error)
	RetrieveRangeSize(ctx context.Context, ID uuid.UUID) (int64, error)
	RetrieveRangeChunks(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunk, error)
	RetrieveRangeChunkReplicas(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunkReplica, error)
	RetrieveRangeReplicas(ctx context.Context, rngID uuid.UUID) ([]*models.RangeReplica, error)

	ReallocateChunks(ctx context.Context, pks []uuid.UUID, newRangeID uuid.UUID) error
	ReallocateChunkReplicas(ctx context.Context, pk []uuid.UUID, newRangeReplicaID uuid.UUID) error
}

type PersistCluster struct {
	Cluster cluster.Cluster
}

// |||| NEW ||||

func (p *PersistCluster) NewRange(ctx context.Context, nodeID int) (*models.Range, error) {
	c := errutil.NewContextCatcher(ctx)
	r := &models.Range{}
	c.Exec(p.Cluster.NewCreate().Model(r).Exec)
	rr := &models.RangeReplica{RangeID: r.ID, NodeID: nodeID}
	c.Exec(p.Cluster.NewCreate().Model(rr).Exec)
	lease := &models.RangeLease{RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(p.Cluster.NewCreate().Model(lease).Exec)
	r.RangeLease = lease
	return r, nil
}

func (p *PersistCluster) NewRangeReplica(ctx context.Context, rangeID uuid.UUID, nodeID int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rangeID, NodeID: nodeID}
	err := p.Cluster.NewCreate().Model(rr).Exec(ctx)
	return rr, err
}

// |||| RETRIEVE ||||

func (p *PersistCluster) RetrieveRange(ctx context.Context, ID uuid.UUID) (*models.Range, error) {
	r := &models.Range{}
	err := p.Cluster.NewRetrieve().Model(r).WherePK(ID).Exec(ctx)
	return r, err
}

func (p *PersistCluster) ccByRangeQ(chunks interface{}, ID uuid.UUID) *cluster.QueryRetrieve {
	return p.Cluster.NewRetrieve().Model(chunks).WhereFields(model.WhereFields{"RangeID": ID})
}

func (p *PersistCluster) RetrieveRangeSize(ctx context.Context, ID uuid.UUID) (int64, error) {
	var size int64 = 0
	err := p.ccByRangeQ(&models.ChannelChunk{}, ID).Calculate(storage.CalculateSum, "Size", &size).Exec(ctx)
	return size, err
}

func (p *PersistCluster) RetrieveRangeChunks(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunk, error) {
	var cc []*models.ChannelChunk
	err := p.ccByRangeQ(cc, rangeID).Exec(ctx)
	return cc, err
}

func (p *PersistCluster) RetrieveRangeReplicas(ctx context.Context, rangeID uuid.UUID) ([]*models.RangeReplica, error) {
	var rr []*models.RangeReplica
	err := p.Cluster.NewRetrieve().Model(&rr).WhereFields(model.WhereFields{"RangeID": rangeID}).Exec(ctx)
	return rr, err
}

func (p *PersistCluster) RetrieveRangeChunkReplicas(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunkReplica, error) {
	var ccr []*models.ChannelChunkReplica
	err := p.Cluster.NewRetrieve().Model(&ccr).WhereFields(model.WhereFields{"ChannelChunk.RangeID": rangeID}).Exec(ctx)
	return ccr, err
}

// |||| RE-ALLOCATE ||||

func (p *PersistCluster) ReallocateChunks(ctx context.Context, pks []uuid.UUID, newRangeID uuid.UUID) error {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: newRangeID})
	}
	return p.Cluster.NewUpdate().Model(&cc).Fields("RangeID").Bulk().Exec(ctx)
}

func (p *PersistCluster) ReallocateChunkReplicas(ctx context.Context, pks []uuid.UUID, newRangeReplicaID uuid.UUID) error {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: newRangeReplicaID})
	}
	return p.Cluster.NewUpdate().Model(&ccr).Fields("RangeReplicaID").Bulk().Exec(ctx)
}
