package config

import (
	"github.com/arya-analytics/aryacore/pkg/ds"
)

type Config struct {
	DS           ds.ConfigChain
	BaseEndpoint []string
}

const (
	AryaDB = "aryadb"
)

func GetConfig() *Config {
	DS := map[string]ds.Config{
		AryaDB: {
			Engine: ds.Postgres,
			Name:   "arya-db",
			Host:   "arya-db",
			Port:   "5628",
			Secure: false,
			Auth: ds.AuthConfig{
				Mode:     ds.Credentials,
				User:     "arya-db-master",
				Password: "arya-dummy-pass",
			},
		},
	}
	return &Config{
		BaseEndpoint: []string{"api"},
		DS:           DS,
	}
}
