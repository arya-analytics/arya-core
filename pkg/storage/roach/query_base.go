package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/uptrace/bun"
)

type queryBase struct {
	exchange *model.Exchange
	catcher  *errutil.Catcher
	handler  storage.ErrorHandler
	db       *bun.DB
}

func (q *queryBase) baseInit(db *bun.DB) {
	q.db = db
	q.catcher = &errutil.Catcher{}
	q.handler = newErrorHandler()
}

func (q *queryBase) baseModel(m interface{}) {
	q.exchange = model.NewExchange(m, catalog().New(m))
}

func (q *queryBase) Dest() *model.Reflect {
	return q.exchange.Dest
}

func (q *queryBase) baseExchangeToSource() {
	q.exchange.ToSource()
}

func (q *queryBase) baseSQL() sqlGen {
	return sqlGen{q.db, q.Dest()}
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
