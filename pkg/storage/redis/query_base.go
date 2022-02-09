package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type queryBase struct {
	_client       *timeseries.Client
	modelExchange *modelExchange
	catcher       *errutil.Catcher
	handler       storage.ErrorHandler
}

func (q *queryBase) baseInit(client *timeseries.Client) {
	q.catcher = &errutil.Catcher{}
	q.handler = newErrorHandler()
	q._client = client
}

func (q *queryBase) baseClient() *timeseries.Client {
	return q._client
}

func (q *queryBase) baseExchangeToDest() {
	q.modelExchange.ToDest()
}

func (q *queryBase) baseExchangeToSource() {
	q.modelExchange.ToSource()
}

func (q *queryBase) baseModel(m interface{}) {
	q.modelExchange = newWrappedModelExchange(storage.NewModelExchange(m,
		catalog().New(m)))
}

func (q *queryBase) baseErr() error {
	return q.handler.Exec(q.catcher.Error())
}
