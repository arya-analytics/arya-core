package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type adapter struct {
	id     uuid.UUID
	client *minio.Client
	driver Driver
}

func newAdapter(driver Driver) (*adapter, error) {
	a := &adapter{
		id:     uuid.New(),
		driver: driver,
	}

	return a, a.open()
}

func bindAdapter(a storage.Adapter) (*adapter, bool) {
	me, ok := a.(*adapter)
	return me, ok
}

func conn(a storage.Adapter) *minio.Client {
	me, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind minio adapter")
	}
	return me.conn()
}

func (a *adapter) ID() uuid.UUID {
	return a.id
}

func (a *adapter) DemandCap() int {
	return a.driver.DemandCap()
}

func (a *adapter) conn() *minio.Client {
	return a.client
}

func (a *adapter) open() error {
	var err error
	a.client, err = a.driver.Connect()
	return err
}
