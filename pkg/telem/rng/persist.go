package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
)

type Persist interface {
	NewRange(ctx context.Context, nodeID int) (*models.Range, error)
}

type PersistCluster struct {
	Cluster cluster.Cluster
}

func (p *PersistCluster) NewRange(ctx context.Context, nodeID int) (*models.Range, error) {
	r := &models.Range{}
	if err := p.Cluster.NewCreate().Model(r).Exec(ctx); err != nil {
		return nil, err
	}
	rr := &models.RangeReplica{
		RangeID: r.ID,
		NodeID:  nodeID,
	}
	if err := p.Cluster.NewCreate().Model(rr).Exec(ctx); err != nil {
		return nil, err
	}
	lease := &models.RangeLease{
		RangeID:        r.ID,
		RangeReplicaID: rr.ID,
		RangeReplica:   rr,
	}
	if err := p.Cluster.NewCreate().Model(lease).Exec(ctx); err != nil {
		return nil, err
	}
	r.RangeLease = lease
	return r, nil
}
