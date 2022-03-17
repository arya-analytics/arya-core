package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/google/uuid"
)

type adapter struct {
	id     uuid.UUID
	client *timeseries.Client
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

func conn(a internal.Adapter) *timeseries.Client {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind redis adapter.")
	}
	return ra.conn()
}

func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) DemandCap() int {
	return a.driver.DemandCap()
}

func (a *adapter) open() error {
	var err error
	a.client, err = a.driver.Connect()
	return err
}

func (a *adapter) conn() *timeseries.Client {
	return a.client
}
