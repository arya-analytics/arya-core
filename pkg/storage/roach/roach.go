package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/uptrace/bun"
)

const (
	Key = "roach"
)

type Engine struct {
	pooler storage.Pooler
}

func (e Engine) conn() *bun.DB {
	a := e.pooler.Retrieve(Key)
	return a.Conn().(*bun.DB)
}

func (e Engine) NewRetrieve() *Retrieve {
	return NewRetrieve(e.conn())
}

func (e Engine) NewCreate() *Create {
	return NewCreate(e.conn())
}












