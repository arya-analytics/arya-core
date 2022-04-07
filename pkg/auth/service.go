package auth

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	qExec query.Execute
}

func NewService(qExec query.Execute) *Service {
	return &Service{qExec: qExec}
}

func (s *Service) Login(ctx context.Context, username, password string) (*models.User, error) {
	user := &models.User{}
	c := errutil.NewCatchSimple(errutil.WithConvert(newErrorConvert()))
	c.Exec(func() error { return s.retrieveUserByUsername(ctx, user, username) })
	c.Exec(func() error { return compareHashAndPassword(user.Password, password) })
	return user, c.Error()
}

func (s *Service) retrieveUserByUsername(ctx context.Context, user *models.User, username string) error {
	return query.NewRetrieve().
		BindExec(s.qExec).
		WhereFields(query.WhereFields{"Username": username}).
		Model(user).
		Exec(ctx)
}
