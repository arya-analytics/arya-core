package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type Adapter struct {
	id     uuid.UUID
	client *minio.Client
	cfg    Config
}

func newAdapter(cfg Config) *Adapter {
	a := &Adapter{
		id:  uuid.New(),
		cfg: cfg,
	}
	a.open()
	return a
}

func bindAdapter(a storage.Adapter) (*Adapter, bool) {
	ma, ok := a.(*Adapter)
	return ma, ok
}

func conn(a storage.Adapter) *minio.Client {
	ma, ok := bindAdapter(a)
	if !ok {
		panic("couldn't bind minio adapter")
	}
	return ma.conn()
}

func (a *Adapter) ID() uuid.UUID {
	return a.id
}

func (a *Adapter) conn() *minio.Client {
	return a.client
}

func (a *Adapter) open() {
	switch a.cfg.Driver {
	case DriverMinIO:
		a.client = connectToMinIO(a.cfg)
	}

}

func connectToMinIO(cfg Config) *minio.Client {
	client, err := minio.New(
		cfg.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.UseTLS,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
