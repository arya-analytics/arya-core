package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"time"
)

type adapter struct {
	client     *timeseries.Client
	driver     Driver
	demand     pool.Demand
	expiration pool.Expire
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		driver:     driver,
		expiration: pool.Expire{Duration: driver.Expiration(), Start: time.Now()},
		demand:     pool.Demand{Max: driver.DemandCap()},
	}
	return a, a.open()
}

func client(a pool.Adapt[*Engine]) *timeseries.Client {
	return a.(*adapter).client
}

func (a *adapter) Match(e *Engine) bool {
	return true
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
