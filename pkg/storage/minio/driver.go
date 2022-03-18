package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"time"
)

type Driver interface {
	Connect() (*minio.Client, error)
	DemandCap() int
	Expiration() time.Duration
}

// |||| CONFIG ||||

type Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Token      string
	UseTLS     bool
	DemandCap  int
	Expiration time.Duration
}

const (
	driverMinioEndpoint   = "endpoint"
	driverMinioAccessKey  = "accessKey"
	driverMinioSecretKey  = "secretKey"
	driverMinioToken      = "token"
	driverMinioUseTLS     = "useTLS"
	driverMinioDemandCap  = "demandCap"
	driverMinioExpiration = "expiration"
)

func configLeaf(name string) string {
	return storage.ConfigTree().Object().Leaf(name)
}

func (c Config) Viper() Config {
	return Config{
		Endpoint:   viper.GetString(configLeaf(driverMinioEndpoint)),
		AccessKey:  viper.GetString(configLeaf(driverMinioAccessKey)),
		SecretKey:  viper.GetString(configLeaf(driverMinioSecretKey)),
		Token:      viper.GetString(configLeaf(driverMinioToken)),
		UseTLS:     viper.GetBool(configLeaf(driverMinioUseTLS)),
		DemandCap:  viper.GetInt(configLeaf(driverMinioDemandCap)),
		Expiration: time.Duration(viper.GetInt(configLeaf(driverMinioExpiration))) * time.Second,
	}

}

// |||| DRIVER ||||

type DriverMinio struct {
	Config Config
}

func (d DriverMinio) Connect() (*minio.Client, error) {
	return minio.New(d.Config.Endpoint, d.buildConfig())
}

func (d DriverMinio) buildConfig() *minio.Options {
	return &minio.Options{
		Creds:  credentials.NewStaticV4(d.Config.AccessKey, d.Config.SecretKey, ""),
		Secure: d.Config.UseTLS,
	}
}

func (d DriverMinio) DemandCap() int {
	return d.Config.DemandCap
}

func (d DriverMinio) Expiration() time.Duration {
	return d.Config.Expiration
}
