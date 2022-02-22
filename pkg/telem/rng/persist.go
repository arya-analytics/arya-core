package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Persist interface {
	NewRange(ctx context.Context, nodeID int) (*models.Range, error)
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
