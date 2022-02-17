package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
)

type Service struct {
	QueryRequest *internal.QueryRequest
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CanHandle(q *internal.QueryRequest) bool {
	return true
}

func (s *Service) Exec(ctx context.Context, q *internal.QueryRequest) error {
	s.QueryRequest = q
	return nil
}
