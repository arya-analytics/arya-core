package aryacore

import (
	"github.com/arya-analytics/aryacore/ds"
	"os"
)

var DS = ds.Configs{
	"default": {
		Engine:   "github.com/uptrace/bun/driver/pgdriver",
		Name:     os.Getenv("ARYA_DB_NAME"),
		Host:     os.Getenv("ARYA_DB_HOST"),
		Port:     os.Getenv("ARYA_DB_PORT"),
		User:     os.Getenv("ARYA_DB_USER"),
		Password: os.Getenv("ARYA_DB_PASSWORD"),
		Secure:   false,
	},
}

func GetConfig() *Config {
	return &Config {
		ds: DS,
	}
}