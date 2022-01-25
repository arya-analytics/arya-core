package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

// |||| QUERY ||||

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
	if err := r.whereBaseValidateReq(); err != nil {
		return r.baseHandleExecErr(err)
	}
	var dvc DataValueChain
	for _, pk := range r.PKs {
		resObj, err := r.baseClient().GetObject(ctx, r.Bucket(), pk.String(), minio.GetObjectOptions{})
		if err != nil {
			return r.baseHandleExecErr(err)
		}
		if vErr := r.validateRes(resObj); vErr != nil {
			return r.baseHandleExecErr(vErr)
		}
		dvc = append(dvc, &DataValue{PK: pk, Data: &Object{resObj}})
	}
	r.baseBindVals(dvc)
	r.baseAdaptToSource()
	return nil
}

func (r *retrieveQuery) validateRes(resObj *minio.Object) error {
	return retrieveQueryResValidator.Exec(resObj)

}

// |||| VALIDATORS ||||

var retrieveQueryResValidator = validate.New([]validate.Func{
	validateResStat,
})

func validateResStat(v interface{}) error {
	res := v.(*minio.Object)
	_, err := res.Stat()
	return err
}
