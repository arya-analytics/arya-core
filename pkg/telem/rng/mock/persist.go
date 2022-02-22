package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
)

type Persist struct {
	Ranges []*models.Range
}

func (p *Persist) NewRange(ctx context.Context, nodeID int) (*models.Range, error) {
	id := uuid.New()
	r := &models.Range{
		ID:     id,
		Status: models.RangeStatusOpen,
		RangeLease: &models.RangeLease{
			ID:      uuid.New(),
			RangeID: id,
			RangeReplica: &models.RangeReplica{
				ID:      uuid.New(),
				RangeID: id,
				NodeID:  nodeID,
			},
		},
	}
	p.Ranges = append(p.Ranges, r)
	return r, nil
}
