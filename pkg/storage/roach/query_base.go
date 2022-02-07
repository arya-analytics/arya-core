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

type queryBase struct {
	exchange *storage.ModelExchange
	catcher  *errutil.Catcher
	handler  storage.ErrorHandler
}

func (q *queryBase) baseInit() {
	q.catcher = &errutil.Catcher{}
	q.handler = newErrorHandler()
}

func (q *queryBase) baseModel(m interface{}) *model.Reflect {
	q.exchange = storage.NewModelExchange(m, catalog().New(m))
	return q.exchange.Dest
}

func (q *queryBase) baseExchangeToSource() {
	q.exchange.ToSource()
}

func (q *queryBase) baseExchangeToDest() {
	q.exchange.ToDest()
}

func (q *queryBase) baseErr() error {
	return q.handler.Exec(q.catcher.Error())
}

func (q *queryBase) baseExec(af errutil.ActionFunc) {
	q.catcher.Exec(af)
}
