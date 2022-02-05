package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Driver interface {
	Connect() (*minio.Client, error)
}

type DriverMinio struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	UseTLS    bool
}

func (d DriverMinio) Connect() (*minio.Client, error) {
	return minio.New(d.Endpoint, d.buildConfig())
}

func (d DriverMinio) buildConfig() *minio.Options {
	return &minio.Options{
		Creds:  credentials.NewStaticV4(d.AccessKey, d.SecretKey, ""),
		Secure: d.UseTLS,
	}
}
