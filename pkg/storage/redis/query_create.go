package redis

import (
	"context"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

type createQuery struct {
	baseQuery
}

func newCreate(client *redistimeseries.Client) *createQuery {
	c := &createQuery{}
	c.baseInit(client)
	return c
}

func (c *createQuery) Model(m interface{}) storage.CacheCreateQuery {
	c.baseModel(m)
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	dRfl := c.modelAdapter.Dest()
	switch dRfl.Type() {
	case reflect.TypeOf(&ChannelSample{}):
		dRfl.ForEach(func(rfl *model.Reflect, i int) {
			cs := rfl.Pointer().(*ChannelSample)
			ccPK := model.NewPK(cs.ChannelConfigID)
			c.catcher.Exec(func() error {
				_, err := c.baseClient().Add(ccPK.String(), cs.Timestamp.Unix(),
					float64(cs.Value))
				return err
			})
		})

	case reflect.TypeOf(&ChannelConfig{}):
		dRfl.ForEach(func(rfl *model.Reflect, i int) {
			cc := rfl.Pointer().(*ChannelConfig)
			c.catcher.Exec(func() error {
				return c.baseClient().CreateKeyWithOptions(dRfl.PKField().String(),
					redistimeseries.CreateOptions{RetentionMSecs: cc.Retention})
			})
		})

	}
	return c.baseErr()
}
