package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	modelAdapter *storage.ModelAdapter
}

func (b *baseQuery) model(m interface{}) interface{} {
	b.modelAdapter = storage.NewModelAdapter(m, roachModelFromStorage(m))
	return b.modelAdapter.DestModel()
}

func (b *baseQuery) adaptToSource() {
	if err := b.modelAdapter.ExchangeToSource(); err != nil {
		log.Fatalln(err)
	}
}

func (b *baseQuery) adaptToDest() {
	if err := b.modelAdapter.ExchangeToDest(); err != nil {
		log.Fatalln(err)
	}
}
