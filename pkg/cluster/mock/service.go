package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
)

type Service struct {
	QueryRequest *cluster.QueryRequest
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CanHandle(q *cluster.QueryRequest) bool {
	return true
}

func (s *Service) Exec(ctx context.Context, q *cluster.QueryRequest) error {
	s.QueryRequest = q
	return nil
}
