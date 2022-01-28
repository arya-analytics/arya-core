package redis

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
)

// |||| CONFIG ||||

type Driver int

const (
	DriverRedisTS Driver = iota
)

type Config struct {
	Host     string
	Port     int
	Driver   Driver
	Password string
	Database int
}

func (c Config) addr() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

// |||| ENGINE ||||

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

func (e *Engine) NewTSRetrieve(a storage.Adapter) storage.CacheTSRetrieveQuery {
	return newTSRetrieve(conn(a))
}

func (e *Engine) NewTSCreate(a storage.Adapter) storage.CacheTSCreateQuery {
	return newTSCreate(conn(a))
}
