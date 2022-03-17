package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// || ADAPTER ||

type adapter struct {
	id     uuid.UUID
	db     *bun.DB
	driver Driver
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		id:     uuid.New(),
		driver: driver,
	}
	return a, a.open()
}

func bindAdapter(a internal.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func conn(a internal.Adapter) *bun.DB {
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

func (a *adapter) DemandCap() int {
	return a.driver.DemandCap()
}

func (a *adapter) conn() *bun.DB {
	return a.db
}

func (a *adapter) open() error {
	var err error
	a.db, err = a.driver.Connect()
	return newErrorConvert().Exec(err)
}

func UnsafeConn(a internal.Adapter) *bun.DB {
	return conn(a)
}
