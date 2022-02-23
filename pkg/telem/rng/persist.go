package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/google/uuid"
)

type Persist interface {
	// |||| RANGE ||||

	NewRange(ctx context.Context, nodeID int) (*models.Range, error)
	CreateRange(ctx context.Context, rng interface{}) error
	CreateRangeLease(ctx context.Context, rngLease interface{}) error
	CreateRangeReplica(ctx context.Context, rngReplica interface{}) error

	RetrieveRange(ctx context.Context, ID uuid.UUID) (*models.Range, error)
	RetrieveRangeSize(ctx context.Context, ID uuid.UUID) (int64, error)
	RetrieveRangeChunks(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunk, error)
	RetrieveRangeChunkReplicas(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunkReplica, error)

	ReallocateChunks(ctx context.Context, pks interface{}, newRangeID uuid.UUID) error
	ReallocateChunkReplicas(ctx context.Context, pks interface{}, newRangeReplicaID uuid.UUID) error

	RetrieveRangeReplicas(ctx context.Context, rngID uuid.UUID) ([]*models.RangeReplica, error)
}

type PersistCluster struct {
	Cluster cluster.Cluster
}

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
