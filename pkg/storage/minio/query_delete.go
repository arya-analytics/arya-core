package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

var removeObjectOpts = minio.RemoveObjectOptions{}

type deleteQuery struct {
	whereBaseQuery
}

func newDelete(client *minio.Client) *deleteQuery {
	d := &deleteQuery{}
	d.baseInit(client)
	return d
}

func (d *deleteQuery) WherePK(pk interface{}) storage.ObjectDeleteQuery {
	d.whereBasePK(pk)
	return d
}

func (d *deleteQuery) WherePKs(pks interface{}) storage.ObjectDeleteQuery {
	d.whereBasePKs(pks)
	return d
}

func (d *deleteQuery) Model(m interface{}) storage.ObjectDeleteQuery {
	d.baseModel(m)
	return d
}

func (d *deleteQuery) Exec(ctx context.Context) error {
	for _, pk := range d.pkChain {
		d.baseExec(func() error {
			return d.baseClient().RemoveObject(
				ctx,
				d.baseBucket(),
				pk.String(),
				removeObjectOpts,
			)
		})
	}
	return d.baseErr()
}
