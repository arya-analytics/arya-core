package minio

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

var getObjectOpts = minio.GetObjectOptions{}

type queryRetrieve struct {
	queryWhereBase
	dvc dataValueChain
}

func newRetrieve(client *minio.Client) *queryRetrieve {
	q := &queryRetrieve{dvc: dataValueChain{}}
	q.baseInit(client)
	return q
}

func (q *queryRetrieve) Model(m interface{}) storage.QueryObjectRetrieve {
	q.baseModel(m)
	return q
}

func (q *queryRetrieve) WherePKs(pks interface{}) storage.QueryObjectRetrieve {
	q.whereBasePKs(pks)
	return q
}

func (q *queryRetrieve) WherePK(pk interface{}) storage.QueryObjectRetrieve {
	q.whereBasePK(pk)
	return q
}

func (q *queryRetrieve) Exec(ctx context.Context) error {
	q.baseValidateReq()
	for _, pk := range q.pkChain {
		var resObj *minio.Object
		q.baseExec(func() (err error) {
			resObj, err = q.baseClient().GetObject(
				ctx,
				q.baseBucket(),
				pk.String(),
				getObjectOpts,
			)
			return err
		})
		q.validateRes(resObj)
		bulk := telem.NewBulk([]byte{})
		q.baseExec(func() error {
			_, err := bulk.ReadFrom(resObj)
			return err

		})
		q.appendToDVC(&dataValue{PK: pk, Data: bulk})
	}
	q.baseBindVals(q.dvc)
	q.baseExchangeToSource()
	return q.baseErr()
}

func (q *queryRetrieve) appendToDVC(dv *dataValue) {
	q.baseExec(func() error {
		q.dvc = append(q.dvc, dv)
		return nil
	})
}

func (q *queryRetrieve) validateRes(resObj *minio.Object) {
	q.baseExec(func() error {
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
