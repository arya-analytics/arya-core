package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type queryBase struct {
	_client       *timeseries.Client
	modelExchange *exchange
	catcher       *errutil.CatchSimple
}

func (q *queryBase) baseInit(client *timeseries.Client) {
	q.catcher = errutil.NewCatchSimple()
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
	q.modelExchange = wrapExchange(model.NewExchange(m,
		catalog().New(m)))
}

func (q *queryBase) baseErr() error {
	return newErrorConvertChain().Exec(q.catcher.Error())
}
