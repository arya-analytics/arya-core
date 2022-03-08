package redis

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Config struct {
	Host      string
	Port      int
	Password  string
	Database  int
	Username  string
	DemandCap int
}

const (
	driverRedisHost      = "host"
	driverRedisPort      = "port"
	driverRedisDatabase  = "database"
	driverRedisPassword  = "password"
	driverRedisUsername  = "username"
	driverRedisDemandCap = "demandCap"
)

func configLeaf(name string) string {
	return storage.ConfigTree().Cache().Leaf(name)
}

func (c Config) Viper() Config {
	return Config{
		Host:      viper.GetString(configLeaf(driverRedisHost)),
		Port:      viper.GetInt(configLeaf(driverRedisPort)),
		Database:  viper.GetInt(configLeaf(driverRedisDatabase)),
		Password:  viper.GetString(configLeaf(driverRedisPassword)),
		Username:  viper.GetString(configLeaf(driverRedisUsername)),
		DemandCap: viper.GetInt(configLeaf(driverRedisDemandCap)),
	}
}

type Driver interface {
	Connect() (*timeseries.Client, error)
	DemandCap() int
}

type DriverRedis struct {
	Config Config
}

func (d DriverRedis) DemandCap() int {
	return d.Config.DemandCap
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
		Username: d.Config.Username,
	}
}
