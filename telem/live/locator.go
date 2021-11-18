package live

import (
	"context"
	"github.com/arya-analytics/aryacore/ds"
	"github.com/arya-analytics/aryacore/telem"
	"github.com/uptrace/bun"
)

type Locator interface {
	getConn(ID int32) ds.Conn
}

func NewLocator(cp *ds.ConnPooler) Locator {
	return &DatabaseLocator{cp}
}

type DatabaseLocator struct {
	cp *ds.ConnPooler
}

func (dl *DatabaseLocator) getConn(ID int32) ds.Conn {
	db := dl.cp.GetOrCreate("aryadb").(*bun.DB)
	ctx := context.Background()
	err := db.NewSelect().Model(new(telem.ChannelConfig)).Where("ID = ",
		ID).Scan(ctx)
	if err != nil {
		panic(err)
	}
	cfg := ds.Config{
		Engine: ds.GorillaWS,
		Name:   "/api/telem/",
		Host:   "10.2.1.22",
	}
	dl.cp.SetConfig("gorillaws", cfg)
	return dl.cp.GetOrCreate("gorwillws")
}
