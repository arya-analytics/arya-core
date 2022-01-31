package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type baseQuery struct {
	_client      *minio.Client
	modelAdapter *storage.ModelAdapter
	catcher      *errutil.Catcher
}

func (b *baseQuery) baseInit(client *minio.Client) {
	b.catcher = &errutil.Catcher{}
	b._client = client
}

func (b *baseQuery) baseClient() *minio.Client {
	return b._client
}

func (b *baseQuery) baseModel(m interface{}) {
	b.modelAdapter = storage.NewModelAdapter(m, catalog().New(m))
}

func (b *baseQuery) baseModelWrapper() *ModelWrapper {
	return &ModelWrapper{rfl: b.modelAdapter.Dest()}
}

func (b *baseQuery) baseBucket() string {
	return b.baseModelWrapper().Bucket()
}

func (b *baseQuery) baseAdaptToSource() {
	b.modelAdapter.ExchangeToSource()
}

func (b *baseQuery) baseAdaptToDest() {
	b.modelAdapter.ExchangeToDest()
}

func (b *baseQuery) baseBindVals(dvc DataValueChain) {
	b.baseModelWrapper().BindDataVals(dvc)
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
	if b.modelAdapter == nil {
		return storage.Error{Type: storage.ErrTypeInvalidArgs,
			Message: "no model provided to query"}
	}
	return nil
}
