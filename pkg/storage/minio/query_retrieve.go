package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/minio/minio-go/v7"
	"reflect"
)

type retrieveQuery struct {
	baseQuery
	pks []model.PK
}

func newRetrieve(client *minio.Client) *retrieveQuery {
	r := &retrieveQuery{}
	r.baseInit(client)
	return r
}

func (r *retrieveQuery) WherePK(pk interface{}) storage.ObjectRetrieveQuery {
	r.pks = append(r.pks, model.NewPK(pk))
	return r
}

func (r *retrieveQuery) WherePKs(pks interface{}) storage.ObjectRetrieveQuery {
	rfl := reflect.ValueOf(pks)
	for i := 0; i < rfl.Len(); i++ {
		r.WherePK(rfl.Index(i).Interface())
	}
	return r
}

func (r *retrieveQuery) Model(m interface{}) storage.ObjectRetrieveQuery {
	r.baseModel(m)
	return r
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	var dvc DataValueChain
	for _, pk := range r.pks {
		obj, err := r.baseClient().GetObject(ctx, r.Bucket(), pk.String(), minio.GetObjectOptions{})
		if err != nil {
			return r.baseHandleExecErr(err)
		}
		dvc = append(dvc, &DataValue{PK: pk, Data: &Object{obj}})
	}
	r.baseBindVals(dvc)
	r.baseAdaptToSource()
	return nil
}
