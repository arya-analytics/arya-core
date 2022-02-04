package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

const (
	pkEqualsSQL  = "ID = ?"
	pkChainInSQL = "ID in (?)"
)

type baseQuery struct {
	exchange *storage.ModelExchange
	catcher  *errutil.Catcher
}

func (b *baseQuery) baseInit() {
	b.catcher = &errutil.Catcher{}
}

func (b *baseQuery) baseModel(m interface{}) *model.Reflect {
	b.exchange = storage.NewModelExchange(m, catalog().New(m))
	return b.exchange.Dest
}

func (b *baseQuery) baseExchangeToSource() {
	b.exchange.ToSource()
}

func (b *baseQuery) baseExchangeToDest() {
	b.exchange.ToDest()
}

func (b *baseQuery) baseErr() error {
	return parseBunErr(b.catcher.Error())
}
