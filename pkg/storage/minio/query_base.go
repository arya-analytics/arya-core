package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	_client      *minio.Client
	modelAdapter *storage.ModelAdapter
	err          error
}

func (b *baseQuery) baseInit(client *minio.Client) {
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

func (b *baseQuery) Bucket() string {
	return b.baseModelWrapper().Bucket()
}

func (b *baseQuery) baseAdaptToSource() {
	if err := b.modelAdapter.ExchangeToSource(); err != nil {
		log.Fatalln(err)
	}
}

func (b *baseQuery) baseAdaptToDest() {
	if err := b.modelAdapter.ExchangeToDest(); err != nil {
		log.Fatalln(err)
	}
}

func (b *baseQuery) baseBindVals(dvc DataValueChain) {
	b.baseModelWrapper().BindDataVals(dvc)
}

func (b *baseQuery) baseBindErr(e error) {
	b.err = e
}

func (b *baseQuery) baseHandleExecErr(e error) error {
	pe := parseMinioErr(e)
	b.baseBindErr(pe)
	return pe
}
