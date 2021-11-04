package config

import (
	"github.com/arya-analytics/arya-core/ds"
	"os"
)


func GetParamsConfig() ds.ConnParamsConfig {
	paramsConfig := ds.ConnParamsConfig {
		"default": {
			Engine:   "github.com/uptrace/bun/driver/pgdriver",
			Name:     os.Getenv("ARYA_DB_NAME"),
			Host:     os.Getenv("ARYA_DB_HOST"),
			Port:     os.Getenv("ARYA_DB_PORT"),
			User:     os.Getenv("ARYA_DB_USER"),
			Password: os.Getenv("ARYA_DB_PASSWORD"),
			Secure: false,
		},
	}
	return paramsConfig
}
