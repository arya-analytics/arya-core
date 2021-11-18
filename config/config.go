package config

import (
	"github.com/arya-analytics/aryacore/ds"
)

type Config struct {
	DS           ds.ConfigChain
	BaseEndpoint []string
}

func GetConfig() *Config {
	DS := map[string]ds.Config{
		"aryadb": {
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
