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
	// CreateRange creates new models.Range, its models.RangeLease, and its lease's models.RangeReplica
	// at the provided leaseNodePK.
	CreateRange(ctx context.Context, leaseNodePK int) (*models.Range, error)
	// CreateRangeReplica creates a new models.RangeReplica on the provided models.Node and belonging to the provided models.Range.
	CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error)
}

type persistRetrieve interface {
	// RetrieveRange retrieves a range by its primary key.
	RetrieveRange(ctx context.Context, pk uuid.UUID) (*models.Range, error)
	// RetrieveRangeSize calculates and returns the size of the range with the provided pk.
	RetrieveRangeSize(ctx context.Context, pk uuid.UUID) (int64, error)
	// RetrieveRangesByStatus retrieves all ranges with the provided models.RangeStatus.
	RetrieveRangesByStatus(ctx context.Context) ([]*models.Range, error)
	// RetrieveRangeChunks retrieves all models.ChannelChunk belonging to the models.Range with primary key rngPK.
	RetrieveRangeChunks(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunk, error)
	// RetrieveRangeChunkReplicas retrieves all models.ChannelChunkReplica belonging to the models.Range
	// with primary key rngPK.
	RetrieveRangeChunkReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunkReplica, error)
	// RetrieveRangeReplicas retrieves al.l models.RangeReplica belonging to the models.Range with primary key rngPK.
	RetrieveRangeReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.RangeReplica, error)
}

type persistUpdate interface {
	// ReallocateChunks reallocates each models.ChannelChunk with a primary key in the slice pks to the models.Range
	// with primary key rngPK.
	ReallocateChunks(ctx context.Context, pks []uuid.UUID, rngPK uuid.UUID) error
	// ReallocateChunkReplicas reallocates each models.ChannelChunkReplica with a primary key in the slice pks to the
	// models.RangeReplica with the primary key RRPK.
	ReallocateChunkReplicas(ctx context.Context, pks []uuid.UUID, RRPK uuid.UUID) error
	// UpdateRangeStatus updates the status of the models.Range with primary key rngPK.
	UpdateRangeStatus(ctx context.Context, rngPK uuid.UUID, status models.RangeStatus) error
}

// Persist persists changes to model made through rng package operations.
type Persist interface {
	persistCreate
	persistRetrieve
	persistUpdate
}

// |||| CLUSTER ||||

// PersistCluster implements Persist and uses a cluster.Cluster as its data store.
type PersistCluster struct {
	Cluster cluster.Cluster
}

// || CREATE ||

func (p *PersistCluster) CreateRange(ctx context.Context, nodePK int) (*models.Range, error) {
	c := errutil.NewContextCatcher(ctx)
	r := &models.Range{Status: models.RangeStatusOpen}
	c.Exec(p.Cluster.NewCreate().Model(r).Exec)
	rr := &models.RangeReplica{RangeID: r.ID, NodeID: nodePK}
	c.Exec(p.Cluster.NewCreate().Model(rr).Exec)
	lease := &models.RangeLease{RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(p.Cluster.NewCreate().Model(lease).Exec)
	r.RangeLease = lease
	return r, c.Error()
}

func (p *PersistCluster) CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rngPK, NodeID: nodePK}
	return rr, p.Cluster.NewCreate().Model(rr).Exec(ctx)
}

// || RETRIEVE ||

func (p *PersistCluster) RetrieveRange(ctx context.Context, pk uuid.UUID) (*models.Range, error) {
	r := &models.Range{}
	return r, p.Cluster.NewRetrieve().
		Model(r).
		WherePK(pk).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID").
		Exec(ctx)
}

func (p *PersistCluster) RetrieveRangesByStatus(ctx context.Context) (ranges []*models.Range, err error) {
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

// || RE-ALLOCATE ||

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

// || UPDATE ||

func (p *PersistCluster) UpdateRangeStatus(ctx context.Context, pk uuid.UUID, status models.RangeStatus) error {
	return p.Cluster.NewUpdate().Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk).Exec(ctx)
}
