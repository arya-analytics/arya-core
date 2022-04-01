package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/uptrace/bun"
	"time"
)

// || ADAPTER ||

type adapter struct {
	db         *bun.DB
	driver     Driver
	demand     pool.Demand
	expiration pool.Expire
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		driver:     driver,
		expiration: pool.Expire{Start: time.Now(), Duration: driver.Expiration()},
		demand:     pool.Demand{Max: driver.DemandCap()},
	}
	return a, a.open()
}

func UnsafeDB(a pool.Adapt[internal.Engine]) *bun.DB {
	return a.(*adapter).db
}

func (a *adapter) Acquire() {
	a.demand.Increment()
}

func (a *adapter) Match(e internal.Engine) bool {
	_, ok := e.(*Engine)
	return ok
}

func (a *adapter) Release() {
	a.demand.Decrement()
}

func (a *adapter) Healthy() bool {
	return !a.expiration.Expired() || !a.demand.Exceeded()
}

func (a *adapter) open() error {
	var err error
	a.db, err = a.driver.Connect()
	return newErrorConvert().Exec(err)
}
