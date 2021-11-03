package arya_core

import "os"

type DatabaseConnParams struct {
	engine   string
	name     string
	host     string
	port     string
	user     string
	password string
}

func getDatabaseConnParams(key string) DatabaseConnParams {
	databases := map[string]DatabaseConnParams{
		"default": {
			engine:   "github.com/uptrace/bun/driver/pgdriver",
			name:     os.Getenv("ARYA_DB_NAME"),
			host:     os.Getenv("ARYA_DB_HOST"),
			port:     os.Getenv("ARYA_DB_PORT"),
			user:     os.Getenv("ARYA_DB_USER"),
			password: os.Getenv("ARYA_DB_PASSWORD"),
		},
	}
	return databases[key]
}
