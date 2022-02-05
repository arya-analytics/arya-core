package redis

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/go-redis/redis/v8"
)

type Driver interface {
	Connect() (*timeseries.Client, error)
}

type DriverRedis struct {
	Host     string
	Port     int
	Password string
	Database int
}

func (d DriverRedis) addr() string {
	return fmt.Sprintf("%s:%v", d.Host, d.Port)
}

func (d DriverRedis) Connect() (*timeseries.Client, error) {
	return timeseries.NewWrap(redis.NewClient(d.buildConfig())), nil

}

func (d DriverRedis) buildConfig() *redis.Options {
	return &redis.Options{
		Addr:     d.addr(),
		DB:       d.Database,
		Password: d.Password,
	}
}
