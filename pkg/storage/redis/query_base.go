package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type baseQuery struct {
	_client      *timeseries.Client
	modelAdapter *modelAdapter
	catcher      *errutil.Catcher
}

func (b *baseQuery) baseInit(client *timeseries.Client) {
	b.catcher = &errutil.Catcher{}
	b._client = client
}

func (b *baseQuery) baseClient() *timeseries.Client {
	return b._client
}

func (b *baseQuery) baseAdaptToDest() {
	b.modelAdapter.ExchangeToDest()
}

func (b *baseQuery) baseAdaptToSource() {
	b.modelAdapter.ExchangeToSource()
}

func (b *baseQuery) baseModel(m interface{}) {
	b.modelAdapter = newWrappedModelAdapter(storage.NewModelAdapter(m,
		catalog().New(m)))
}

func (b *baseQuery) baseErr() error {
	return parseRedisTSErr(b.catcher.Error())
}
