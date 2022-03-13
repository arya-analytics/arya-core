package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
	// RetrieveOpenRanges retrieves all open ranges.
	RetrieveOpenRanges(ctx context.Context) ([]*models.Range, error)
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

// PersistCluster implements Persist and uses a cluster.clust as its data store.
type PersistCluster struct {
	clust cluster.Cluster
}

func NewPersistCluster(clust cluster.Cluster) *PersistCluster {
	return &PersistCluster{clust: clust}
}

// || CREATE ||

func createRange(ctx context.Context, exec query.Execute, nodePK int) (*models.Range, error) {
	r := &models.Range{Status: models.RangeStatusOpen}
	c := errutil.NewCatchContext(ctx)
	c.Exec(query.NewCreate().BindExec(exec).Model(r).Exec)
	rr := &models.RangeReplica{RangeID: r.ID, NodeID: nodePK}
	c.Exec(query.NewCreate().BindExec(exec).Model(rr).Exec)
	lease := &models.RangeLease{RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(query.NewCreate().BindExec(exec).Model(lease).Exec)
	r.RangeLease = lease
	return r, c.Error()
}

func createRangeReplica(ctx context.Context, qExec query.Execute, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rngPK, NodeID: nodePK}
	return rr, query.NewCreate().Model(rr).BindExec(qExec).Exec(ctx)
}

// || RETRIEVE ||

func retrieveRangeQuery(qExec query.Execute, rng *models.Range, pk uuid.UUID) *query.Retrieve {
	return rangeLeaseReplicaRelationQuery(qExec, rng).WherePK(pk)
}

func retrieveRangeReplicasQuery(qExec query.Execute, rr []*models.RangeReplica, rngPK uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().BindExec(qExec).Model(&rr).WhereFields(query.WhereFields{"RangeID": rngPK})
}

func openRangeQuery(qExec query.Execute, rng []*models.Range) *query.Retrieve {
	return rangeLeaseReplicaRelationQuery(qExec, &rng).WhereFields(query.WhereFields{"Status": models.RangeStatusOpen})
}

func rangeLeaseReplicaRelationQuery(qExec query.Execute, rng interface{}) *query.Retrieve {
	return query.NewRetrieve().
		BindExec(qExec).
		Model(rng).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID")
}

func retrieveRangeChunksQuery(qExec query.Execute, chunks interface{}, pk uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().BindExec(qExec).Model(chunks).WhereFields(query.WhereFields{"RangeID": pk})
}

func retrieveRangeSizeQuery(qExec query.Execute, pk uuid.UUID, into *int64) *query.Retrieve {
	return retrieveRangeChunksQuery(qExec, &models.ChannelChunk{}, pk).Calc(query.CalcSum, "Size", &into)
}

func updateRangeStatusQuery(qExec query.Execute, pk uuid.UUID, status models.RangeStatus) *query.Update {
	return query.NewUpdate().BindExec(qExec).Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk)
}

func reallocateChunksQuery(qExec query.Execute, pks []uuid.UUID, rngPK uuid.UUID) *query.Update {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: rngPK})
	}
	return query.NewUpdate().BindExec(qExec).Model(&cc).Fields("RangeID").Bulk()
}

func retrieveRangeChunkReplicasQuery(qExec query.Execute, ccr []*models.ChannelChunkReplica, rngPK uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().
		BindExec(qExec).
		Model(&ccr).
		WhereFields(query.WhereFields{"ChannelChunk.RangeID": rngPK}).
		Fields("ID", "RangeReplicaID")
}

// || RE-ALLOCATE ||

func reallocateChunkReplicasQuery(qExec query.Execute, pks []uuid.UUID, rrPK uuid.UUID) *query.Update {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: rrPK})
	}
	return query.NewUpdate().BindExec(qExec).Model(&ccr).Fields("RangeReplicaID").Bulk()
}
