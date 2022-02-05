package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

// || ADAPTER ||

type adapter struct {
	id     uuid.UUID
	db     *bun.DB
	driver Driver
}

func newAdapter(driver Driver) *adapter {
	a := &adapter{
		id:     uuid.New(),
		driver: driver,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func conn(a storage.Adapter) *bun.DB {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind roach adapter.")
	}
	return ra.conn()
}

// ID implements the storage.Adapter interface.
func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) conn() *bun.DB {
	return a.db
}

func (a *adapter) open() {
	var err error
	a.db, err = a.driver.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}
