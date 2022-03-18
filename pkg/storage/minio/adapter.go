package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"time"
)

type adapter struct {
	id         uuid.UUID
	client     *minio.Client
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
	me, ok := a.(*adapter)
	return me, ok
}

func conn(a internal.Adapter) *minio.Client {
	ma, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind minio adapter")
	}
	return ma.client
}

func (a *adapter) Acquire() {
	a.demand.Increment()
}

func (a *adapter) Release() {
	a.demand.Decrement()
}

func (a *adapter) Healthy() bool {
	return !a.expiration.Expired() && !a.demand.Exceeded()
}

func (a *adapter) open() error {
	var err error
	a.client, err = a.driver.Connect()
	return err
}
