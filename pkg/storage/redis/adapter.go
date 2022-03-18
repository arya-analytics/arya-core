package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/google/uuid"
	"time"
)

type adapter struct {
	id         uuid.UUID
	client     *timeseries.Client
	driver     Driver
	demand     internal.Demand
	expiration internal.Expiration
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		id:     uuid.New(),
		driver: driver,
		expiration: internal.Expiration{
			Duration: driver.Expiration(),
			Start:    time.Now(),
		},
		demand: internal.Demand{Max: driver.DemandCap()},
	}
	return a, a.open()
}

func bindAdapter(a internal.Adapter) (*adapter, bool) {
	ra, ok := a.(*adapter)
	return ra, ok
}

func conn(a internal.Adapter) *timeseries.Client {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind redis adapter.")
	}
	return ra.client
}

func (a *adapter) Acquire() {
	a.demand.Increment()
}

func (a *adapter) Release() {
	a.demand.Decrement()
}

func (a *adapter) Healthy() bool {
	return !a.expiration.Expired() || a.demand.Exceeded()
}

func (a *adapter) open() error {
	var err error
	a.client, err = a.driver.Connect()
	return err
}
