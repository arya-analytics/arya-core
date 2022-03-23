package roach

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// || ADAPTER ||

type adapter struct {
	id         uuid.UUID
	db         *bun.DB
	driver     Driver
	demand     internal.Demand
	expiration internal.Expiration
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		id:     uuid.New(),
		driver: driver,
		expiration: internal.Expiration{
			Start:    time.Now(),
			Duration: driver.Expiration(),
		},
		demand: internal.Demand{Max: driver.DemandCap()},
	}
	return a, a.open()
}

func UnsafeDB(a internal.Adapter) *bun.DB {
	return a.(*adapter).db
}

func (a *adapter) Acquire() {
	a.demand.Increment()
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
