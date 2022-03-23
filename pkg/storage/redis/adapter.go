package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"time"
)

type adapter struct {
	client     *timeseries.Client
	driver     Driver
	demand     internal.Demand
	expiration internal.Expiration
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		driver:     driver,
		expiration: internal.Expiration{Duration: driver.Expiration(), Start: time.Now()},
		demand:     internal.Demand{Max: driver.DemandCap()},
	}
	return a, a.open()
}

func client(a internal.Adapter) *timeseries.Client {
	return a.(*adapter).client
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
