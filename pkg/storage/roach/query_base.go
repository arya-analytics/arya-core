package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	modelAdapter *storage.SingleModelAdapter
}

func (b *baseQuery) baseModel(m interface{}) interface{} {
	opts := &storage.ModelAdapterOpts{
		Source: m,
		Dest:   newRoachModelFromStorage(m),
	}
	b.modelAdapter = storage.NewSingleModelAdapter(opts)
	return b.modelAdapter.Dest()
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
