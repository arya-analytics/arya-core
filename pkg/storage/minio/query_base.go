package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type queryBase struct {
	_client       *minio.Client
	modelExchange *modelExchange
	catcher       *errutil.Catcher
}

func (q *queryBase) baseInit(client *minio.Client) {
	q.catcher = &errutil.Catcher{}
	q._client = client
}

func (q *queryBase) baseClient() *minio.Client {
	return q._client
}

func (q *queryBase) baseModel(m interface{}) {
	q.modelExchange = newWrappedModelExchange(storage.NewModelExchange(m,
		catalog().New(m)))
}

func (q *queryBase) baseBucket() string {
	return q.modelExchange.Bucket()
}

func (q *queryBase) baseExchangeToSource() {
	q.modelExchange.ToSource()
}

func (q *queryBase) baseExchangeToDest() {
	q.modelExchange.ToDest()
}

func (q *queryBase) baseBindVals(dvc dataValueChain) {
	q.modelExchange.BindDataVals(dvc)
}

func (q *queryBase) baseExec(af errutil.ActionFunc) {
	q.catcher.Exec(af)
}

func (q *queryBase) baseErr() error {
	return parseMinioErr(q.catcher.Error())
}

func (q *queryBase) baseValidateReq() {
	q.catcher.Exec(func() error { return baseQueryReqValidator.Exec(q) })
}

// |||| VALIDATORS |||

var baseQueryReqValidator = validate.New([]validate.Func{
	validateModelProvided,
})

func validateModelProvided(v interface{}) error {
	b := v.(*queryBase)
	if b.modelExchange == nil {
		return storage.Error{Type: storage.ErrTypeInvalidArgs,
			Message: "no model provided to query"}
	}
	return nil
}
