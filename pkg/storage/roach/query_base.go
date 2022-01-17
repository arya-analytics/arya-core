package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	modelAdapter storage.ModelAdapter
}

func (b *baseQuery) baseModel(m interface{}) interface{} {
	var err error
	b.modelAdapter, err = storage.NewModelAdapter(m, catalog().New(m))
	if err != nil {
		log.Fatalln(err)
	}
	return b.modelAdapter.DestPointer()
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
