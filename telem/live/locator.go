package live

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/config"
	"github.com/arya-analytics/aryacore/ds"
	"github.com/arya-analytics/aryacore/telem"
	"github.com/uptrace/bun"
)

type Locator interface {
	Locate(ChanCfgIds []int32) []Receiver
}

func NewLocator(pooler *ds.ConnPooler) Locator {
	return &DatabaseLocator{pooler}
}

type DatabaseLocator struct {
	pooler *ds.ConnPooler
}

func (dl DatabaseLocator) Locate(ChanCfgIds []int32) []Receiver {
	ctx := context.Background()
	db := dl.pooler.GetOrCreate(config.AryaDB).(*bun.DB)
	var cfgs []telem.ChannelConfig
	db.NewSelect().Model(&cfgs).Scan(ctx)
	fmt.Println(cfgs)
	return []Receiver{}
}
