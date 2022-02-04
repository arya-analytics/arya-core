package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

var getObjectOpts = minio.GetObjectOptions{}

type retrieveQuery struct {
	whereBaseQuery
	dvc dataValueChain
}

func newRetrieve(client *minio.Client) *retrieveQuery {
	r := &retrieveQuery{dvc: dataValueChain{}}
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
	for _, pk := range r.pkChain {
		var resObj *minio.Object
		r.baseExec(func() (err error) {
			resObj, err = r.baseClient().GetObject(
				ctx,
				r.baseBucket(),
				pk.String(),
				getObjectOpts,
			)
			return err
		})
		r.validateRes(resObj)
		r.appendToDVC(&dataValue{PK: pk, Data: &object{resObj}})
	}
	r.baseBindVals(r.dvc)
	r.baseExchangeToSource()
	return r.baseErr()
}

func (r *retrieveQuery) appendToDVC(dv *dataValue) {
	r.baseExec(func() error {
		r.dvc = append(r.dvc, dv)
		return nil
	})
}

func (r *retrieveQuery) validateRes(resObj *minio.Object) {
	r.baseExec(func() error {
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
