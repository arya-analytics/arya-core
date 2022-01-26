package redists

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
)

type Driver int

const (
	DriverRedisTS Driver = iota
)

type Config struct {
	Host     string
	Port     int
	Driver   Driver
	Database string
	Password string
}

func (c Config) addr() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

func (c Config) authString() *string {
	if c.Password != "" {
		authString := fmt.Sprintf("AUTH %s", c.Password)
		return &authString
	}
	return nil
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
