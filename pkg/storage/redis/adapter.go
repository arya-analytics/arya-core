package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	id     uuid.UUID
	client *timeseries.Client
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

func conn(a storage.Adapter) *timeseries.Client {
	ra, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind redis adapter.")
	}
	return ra.conn()
}

func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) open() {
	var err error
	a.client, err = a.driver.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}

func (a *adapter) conn() *timeseries.Client {
	return a.client
}
