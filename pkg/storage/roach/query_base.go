package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	log "github.com/sirupsen/logrus"
)

type baseQuery struct {
	modelAdapter storage.ModelAdapter
	err          error
}

func (b *baseQuery) baseModel(m interface{}) *model.Reflect {
	b.modelAdapter = storage.NewModelAdapter(m, catalog().New(m))
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
	pe := parseBunErr(e)
	b.baseBindErr(pe)
	return pe
}
