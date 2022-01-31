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
	dvc DataValueChain
}

func newRetrieve(client *minio.Client) *retrieveQuery {
	r := &retrieveQuery{dvc: DataValueChain{}}
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
	r.whereBaseValidateReq()
	for _, pk := range r.PKs {
		var resObj *minio.Object
		r.catcher.Exec(func() (err error) {
			resObj, err = r.baseClient().GetObject(ctx, r.baseBucket(), pk.String(),
				minio.GetObjectOptions{})
			return err
		})
		r.validateRes(resObj)
		r.appendToDVC(&DataValue{PK: pk, Data: &Object{resObj}})
	}
	r.baseBindVals(r.dvc)
	r.baseAdaptToSource()
	return r.baseErr()
}

func (r *retrieveQuery) appendToDVC(dv *DataValue) {
	r.catcher.Exec(func() error {
		r.dvc = append(r.dvc, dv)
		return nil
	})
}

func (r *retrieveQuery) validateRes(resObj *minio.Object) {
	r.catcher.Exec(func() error {
		return retrieveQueryResValidator.Exec(resObj)
	})

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
