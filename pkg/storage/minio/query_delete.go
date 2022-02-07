package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var removeObjectOpts = minio.RemoveObjectOptions{}

type queryDelete struct {
	queryWhereBase
}

func newDelete(client *minio.Client) *queryDelete {
	q := &queryDelete{}
	q.baseInit(client)
	return q
}

func (q *queryDelete) WherePK(pk interface{}) storage.QueryObjectDelete {
	q.whereBasePK(pk)
	return q
}

func (q *queryDelete) WherePKs(pks interface{}) storage.QueryObjectDelete {
	q.whereBasePKs(pks)
	return q
}

func (q *queryDelete) Model(m interface{}) storage.QueryObjectDelete {
	q.baseModel(m)
	return q
}

func (q *queryDelete) Exec(ctx context.Context) error {
	for _, pk := range q.pkChain {
		q.baseExec(func() error {
			return q.baseClient().RemoveObject(
				ctx,
				q.baseBucket(),
				pk.String(),
				removeObjectOpts,
			)
		})
	}
	return q.baseErr()
}
