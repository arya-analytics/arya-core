package redis

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     int
	Password string
	Database int
}

const (
	driverRedisHost     = "host"
	driverRedisPort     = "port"
	driverRedisDatabase = "database"
	driverRedisPassword = "password"
)

func configLeaf(name string) string {
	return storage.ConfigTree().Cache().Leaf(name)
}

func (c Config) Viper() Config {
	return Config{
		Host:     viper.GetString(configLeaf(driverRedisHost)),
		Port:     viper.GetInt(configLeaf(driverRedisPort)),
		Database: viper.GetInt(configLeaf(driverRedisDatabase)),
		Password: viper.GetString(configLeaf(driverRedisPassword)),
	}
}

type Driver interface {
	Connect() (*timeseries.Client, error)
}

type DriverRedis struct {
	Config Config
}

func (d DriverRedis) addr() string {
	return fmt.Sprintf("%s:%v", d.Config.Host, d.Config.Port)
}

func (d DriverRedis) Connect() (*timeseries.Client, error) {
	return timeseries.NewWrap(redis.NewClient(d.buildConfig())), nil

}

func (d DriverRedis) buildConfig() *redis.Options {
	return &redis.Options{
		Addr:     d.addr(),
		DB:       d.Config.Database,
		Password: d.Config.Password,
	}
}
