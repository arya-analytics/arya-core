package aryacore

import (
	"github.com/arya-analytics/aryacore/ds"
)

func GetConfig() Config {
	//DS := map[string]ds.Config {
	//	"default": {
	//		Engine:   ds.Postgres,
	//		Name:     os.Getenv("ARYA_DB_NAME"),
	//		Host:     os.Getenv("ARYA_DB_HOST"),
	//		Port:     os.Getenv("ARYA_DB_PORT"),
	//		Secure:   false,
	//		Auth: ds.AuthConfig{
	//			Mode: ds.Credentials,
	//			User: os.Getenv("ARYA_DB_USER"),
	//			Password: os.Getenv("ARYA_DB_PASSWORD"),
	//		},
	//	},
	//}
	DS := map[string]ds.Config {
		"default": {
			Engine:   ds.Postgres,
			Name:     "arya-db",
			Host:    "arya-db" ,
			Port:     "5628",
			Secure:   false,
			Auth: ds.AuthConfig{
				Mode: ds.Credentials,
				User: "arya-db-master",
				Password: "arya-dummy-pass",
			},
		},
	}
	return  Config{
		DS: DS,
	}
}