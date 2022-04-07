package auth

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	qExec query.Execute
}

func NewService(qExec query.Execute) *Service {
	return &Service{qExec: qExec}
}

func (s *Service) Login(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.retrieveUserByUsername(ctx, username)
	if err != nil {
		return nil, newErrorConvert().Exec(err)
	}
	return user, newErrorConvert().Exec(compareHashAndPassword(user.Password, password))
}

func (s *Service) retrieveUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	return user, query.
		NewRetrieve().
		BindExec(s.qExec).
		WhereFields(query.WhereFields{"Username": username}).
		Model(user).
		Exec(ctx)
}
