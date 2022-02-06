package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type baseQuery struct {
	_client       *minio.Client
	modelExchange *modelExchange
	catcher       *errutil.Catcher
}

func (b *baseQuery) baseInit(client *minio.Client) {
	b.catcher = &errutil.Catcher{}
	b._client = client
}

func (b *baseQuery) baseClient() *minio.Client {
	return b._client
}

func (b *baseQuery) baseModel(m interface{}) {
	b.modelExchange = newWrappedModelExchange(storage.NewModelExchange(m,
		catalog().New(m)))
}

func (b *baseQuery) baseBucket() string {
	return b.modelExchange.Bucket()
}

func (b *baseQuery) baseExchangeToSource() {
	b.modelExchange.ToSource()
}

func (b *baseQuery) baseExchangeToDest() {
	b.modelExchange.ToDest()
}

func (b *baseQuery) baseBindVals(dvc dataValueChain) {
	b.modelExchange.BindDataVals(dvc)
}

func (b *baseQuery) baseExec(af errutil.ActionFunc) {
	b.catcher.Exec(af)
}

func (b *baseQuery) baseErr() error {
	return parseMinioErr(b.catcher.Error())
}

func (b *baseQuery) baseValidateReq() {
	b.catcher.Exec(func() error { return baseQueryReqValidator.Exec(b) })
}

// |||| VALIDATORS |||

var baseQueryReqValidator = validate.New([]validate.Func{
	validateModelProvided,
})

func validateModelProvided(v interface{}) error {
	b := v.(*baseQuery)
	if b.modelExchange == nil {
		return storage.Error{Type: storage.ErrTypeInvalidArgs,
			Message: "no model provided to query"}
	}
	return nil
}
