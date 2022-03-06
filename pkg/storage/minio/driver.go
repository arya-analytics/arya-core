package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

// |||| CONFIG ||||

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	UseTLS    bool
}

const (
	driverMinioEndpoint  = "endpoint"
	driverMinioAccessKey = "accessKey"
	driverMinioSecretKey = "secretKey"
	driverMinioToken     = "token"
	driverMinioUseTLS    = "useTLS"
)

func configLeaf(name string) string {
	return storage.ConfigTree().Object().Leaf(name)
}

func (c Config) Viper() Config {
	return Config{
		Endpoint:  viper.GetString(configLeaf(driverMinioEndpoint)),
		AccessKey: viper.GetString(configLeaf(driverMinioAccessKey)),
		SecretKey: viper.GetString(configLeaf(driverMinioSecretKey)),
		Token:     viper.GetString(configLeaf(driverMinioToken)),
		UseTLS:    viper.GetBool(configLeaf(driverMinioUseTLS)),
	}

}

// |||| DRIVER ||||

type Driver interface {
	Connect() (*minio.Client, error)
}

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
