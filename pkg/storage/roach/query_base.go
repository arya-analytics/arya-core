package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type baseQuery struct {
	modelAdapter *storage.ModelAdapter
	catcher      *errutil.Catcher
}

func (b *baseQuery) baseInit() {
	b.catcher = &errutil.Catcher{}
}

func (b *baseQuery) baseModel(m interface{}) *model.Reflect {
	b.modelAdapter = storage.NewModelAdapter(m, catalog().New(m))
	return b.modelAdapter.Dest()
}

func (b *baseQuery) baseAdaptToSource() {
	b.catcher.Exec(b.modelAdapter.ExchangeToSource)
}

func (b *baseQuery) baseAdaptToDest() {
	b.catcher.Exec(b.modelAdapter.ExchangeToDest)
}

func (b *baseQuery) baseErr() error {
	return parseBunErr(b.catcher.Error())
}
