package minio

import "github.com/arya-analytics/aryacore/pkg/storage"

type Driver int

const (
	DriverMinIO Driver = iota
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	UseTLS    bool
	Driver    Driver
}

type Engine struct {
	cfg Config
}

func New(cfg Config) *Engine {
	return &Engine{cfg}
}

func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.cfg)
}

func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) NewCreate(a storage.Adapter) storage.ObjectCreateQuery {
	return newCreate(conn(a))
}

func (e *Engine) NewMigrate(a storage.Adapter) storage.ObjectMigrateQuery {
	return newMigrate(conn(a))
}
