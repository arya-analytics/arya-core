package redis

import (
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type baseQuery struct {
	_client      *redistimeseries.Client
	modelAdapter *storage.ModelAdapter
	catcher      *errutil.Catcher
}

func (b *baseQuery) baseInit(client *redistimeseries.Client) {
	b.catcher = &errutil.Catcher{}
	b._client = client
}

func (b *baseQuery) baseClient() *redistimeseries.Client {
	return b._client
}

func (b *baseQuery) baseAdaptToDest() {
	b.catcher.Exec(b.modelAdapter.ExchangeToDest)
}

func (b *baseQuery) baseAdaptToSource() {
	b.catcher.Exec(b.modelAdapter.ExchangeToSource)
}

func (b *baseQuery) baseModel(m interface{}) {
	b.modelAdapter = storage.NewModelAdapter(m, catalog().New(m))
}

func (b *baseQuery) baseErr() error {
	return parseRedisTSErr(b.catcher.Error())
}
