package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

// |||| SERVICE |||

type ServiceRemote interface {
}

type Service struct {
	qa     query.Assemble
	remote ServiceRemote
}

func NewService(local query.Assemble, remote ServiceRemote) *Service {
	return &Service{remote: remote, qa: local}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	return catalog().Contains(p.Model())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{})
}

func (s *Service) create() {

}

func (s *Service) retrieve() {

}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelSample{}}
}
