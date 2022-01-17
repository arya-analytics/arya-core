package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	modelAdapter storage.ModelAdapter
	err          error
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

func (b *baseQuery) baseBindErr(e error) {
	b.err = e
}

func (b *baseQuery) baseCheckErr() bool {
	return b.err != nil
}

func (b *baseQuery) baseErr() error {
	return b.err
}

func (b *baseQuery) baseHandleExecErr(e error) error {
	if e != nil {
		pe := parseBunErr(e)
		b.baseBindErr(pe)
		return pe
	}
	return nil
}
