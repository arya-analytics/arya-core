package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
)

type retrieveQuery struct {
	whereBaseQuery
}

func newRetrieve(client *minio.Client) *retrieveQuery {
	r := &retrieveQuery{}
	r.baseInit(client)
	return r
}

func (r *retrieveQuery) Model(m interface{}) storage.ObjectRetrieveQuery {
	r.baseModel(m)
	return r
}

func (r *retrieveQuery) WherePKs(pks interface{}) storage.ObjectRetrieveQuery {
	r.whereBasePKs(pks)
	return r
}

func (r *retrieveQuery) WherePK(pk interface{}) storage.ObjectRetrieveQuery {
	r.whereBasePK(pk)
	return r
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	var dvc DataValueChain
	for _, pk := range r.pks {
		obj, err := r.baseClient().GetObject(ctx, r.Bucket(), pk.String(), minio.GetObjectOptions{})
		if err != nil {
			return r.baseHandleExecErr(err)
		}
		stat, err := obj.Stat()
		if stat.Key == "" {
			if len(r.pks) == 1 {
				return storage.NewError(storage.ErrTypeItemNotFound)
			}
		} else {
			dvc = append(dvc, &DataValue{PK: pk, Data: &Object{obj}})
		}
	}
	r.baseBindVals(dvc)
	r.baseAdaptToSource()
	return nil
}
