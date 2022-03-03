package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	QueryRequest *query.Pack
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CanHandle(q *query.Pack) bool {
	return true
}

func (s *Service) Exec(ctx context.Context, q *query.Pack) error {
	s.QueryRequest = q
	return nil
}
