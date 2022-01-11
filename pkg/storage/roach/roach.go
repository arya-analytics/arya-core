package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

const (
	engineRole = storage.EngineRoleMetaData
	engineType = storage.EngineTypeRoach
)

type Engine struct {
	pooler storage.Pooler
}

func NewEngine(pooler storage.Pooler) *Engine {
	return &Engine{pooler: pooler}

}

func (e Engine) Role() storage.EngineRole {
	return engineRole
}

func (e Engine) Type() storage.EngineType {
	return engineType
}

func (e Engine) conn() *bun.DB {
	a, err := e.pooler.Retrieve(storage.EngineTypeRoach)
	if err != nil {
		log.Fatalln(err)
	}
	return a.Conn().(*bun.DB)
}



func (e Engine) NewRetrieve() *Retrieve {
	return NewRetrieve(e.conn())
}

func (e Engine) NewCreate() *Create {
	return NewCreate(e.conn())
}
