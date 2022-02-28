package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type queryBase struct {
	_client       *minio.Client
	modelExchange *modelExchange
	catcher       *errutil.CatchSimple
	handler       storage.ErrorHandler
}

func (q *queryBase) baseInit(client *minio.Client) {
	q.catcher = &errutil.CatchSimple{}
	q.handler = newErrorHandler()
	q._client = client
}

func (q *queryBase) baseClient() *minio.Client {
	return q._client
}

func (q *queryBase) baseModel(m interface{}) {
	q.modelExchange = newWrappedModelExchange(model.NewExchange(m,
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
	return q.handler.Exec(q.catcher.Error())
}

func (q *queryBase) baseValidateReq() {
	q.catcher.Exec(func() error { return baseQueryReqValidator.Exec(q).Error() })
}

// |||| VALIDATORS |||

var baseQueryReqValidator = validate.New[*queryBase]([]func(q *queryBase) error{
	validateModelProvided,
})

func validateModelProvided(q *queryBase) error {
	if q.modelExchange == nil {
		return storage.Error{Type: storage.ErrorTypeInvalidArgs,
			Message: "no model provided to query"}
	}
	return nil
}
