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

type persistCreate interface {
	CreateRange(ctx context.Context, nodePK int) (*models.Range, error)
	CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error)
}

type persistRetrieve interface {
	RetrieveRange(ctx context.Context, pk uuid.UUID) (*models.Range, error)
	RetrieveRangeSize(ctx context.Context, pk uuid.UUID) (int64, error)
	RetrieveOpenRanges(ctx context.Context) ([]*models.Range, error)
	RetrieveRangeChunks(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunk, error)
	RetrieveRangeChunkReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunkReplica, error)
	RetrieveRangeReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.RangeReplica, error)
}

type persistUpdate interface {
	ReallocateChunks(ctx context.Context, pks []uuid.UUID, newRngPK uuid.UUID) error
	ReallocateChunkReplicas(ctx context.Context, pk []uuid.UUID, newRRPK uuid.UUID) error
	UpdateRangeStatus(ctx context.Context, rngPK uuid.UUID, status models.RangeStatus) error
}

type Persist interface {
	persistCreate
	persistRetrieve
	persistUpdate
}

type PersistCluster struct {
	Cluster cluster.Cluster
}

// |||| NEW ||||

func (p *PersistCluster) CreateRange(ctx context.Context, nodePK int) (*models.Range, error) {
	c := errutil.NewContextCatcher(ctx)
	r := &models.Range{Status: models.RangeStatusOpen}
	c.Exec(p.Cluster.NewCreate().Model(r).Exec)
	rr := &models.RangeReplica{RangeID: r.ID, NodeID: nodePK}
	c.Exec(p.Cluster.NewCreate().Model(rr).Exec)
	lease := &models.RangeLease{RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(p.Cluster.NewCreate().Model(lease).Exec)
	r.RangeLease = lease
	return r, nil
}

func (p *PersistCluster) CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rngPK, NodeID: nodePK}
	err := p.Cluster.NewCreate().Model(rr).Exec(ctx)
	return rr, err
}

// |||| RETRIEVE ||||

func (p *PersistCluster) RetrieveRange(ctx context.Context, pk uuid.UUID) (*models.Range, error) {
	r := &models.Range{}
	return r, p.Cluster.NewRetrieve().
		Model(r).
		WherePK(pk).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID").
		Exec(ctx)
}

func (p *PersistCluster) RetrieveOpenRanges(ctx context.Context) (ranges []*models.Range, err error) {
	err = p.Cluster.NewRetrieve().
		Model(&ranges).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID").
		WhereFields(model.WhereFields{"Status": models.RangeStatusOpen}).Exec(ctx)
	return ranges, err

}

func (p *PersistCluster) ccByRangeQ(chunks interface{}, pk uuid.UUID) *cluster.QueryRetrieve {
	return p.Cluster.NewRetrieve().Model(chunks).WhereFields(model.WhereFields{"RangeID": pk})
}

func (p *PersistCluster) RetrieveRangeSize(ctx context.Context, pk uuid.UUID) (int64, error) {
	var size int64 = 0
	return size, p.ccByRangeQ(&models.ChannelChunk{}, pk).Calculate(storage.CalculateSum, "Size", &size).Exec(ctx)
}

func (p *PersistCluster) RetrieveRangeChunks(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunk, error) {
	var cc []*models.ChannelChunk
	return cc, p.ccByRangeQ(&cc, rngPK).Exec(ctx)
}

func (p *PersistCluster) RetrieveRangeReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.RangeReplica, error) {
	var rr []*models.RangeReplica
	return rr, p.Cluster.NewRetrieve().Model(&rr).WhereFields(model.WhereFields{"RangeID": rngPK}).Exec(ctx)
}

func (p *PersistCluster) RetrieveRangeChunkReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunkReplica, error) {
	var ccr []*models.ChannelChunkReplica
	return ccr, p.Cluster.
		NewRetrieve().
		Model(&ccr).
		WhereFields(model.WhereFields{"ChannelChunk.RangeID": rngPK}).
		Fields("ID", "RangeReplicaID").Exec(ctx)
}

// |||| RE-ALLOCATE ||||

func (p *PersistCluster) ReallocateChunks(ctx context.Context, pks []uuid.UUID, newRngPK uuid.UUID) error {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: newRngPK})
	}
	return p.Cluster.NewUpdate().Model(&cc).Fields("RangeID").Bulk().Exec(ctx)
}

func (p *PersistCluster) ReallocateChunkReplicas(ctx context.Context, pks []uuid.UUID, newRRPK uuid.UUID) error {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: newRRPK})
	}
	return p.Cluster.NewUpdate().Model(&ccr).Fields("RangeReplicaID").Bulk().Exec(ctx)
}

// ||| UPDATE |||

func (p *PersistCluster) UpdateRangeStatus(ctx context.Context, pk uuid.UUID, status models.RangeStatus) error {
	return p.Cluster.NewUpdate().Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk).Exec(ctx)
}
