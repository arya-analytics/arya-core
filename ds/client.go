package db

import (
	config "github.com/arya-analytics/arya-core/config"
)
// || SUPPORTED CONNECTION TYPES ||
type DatabaseConnectionUtil

// DBConnManager || CONNECTION MANAGER ||
type DBConnManager struct {
	getConnParams func(key string) config.DatabaseConnParams
}

func (dbcm *DBConnManager) getDSN(key string) string {
	connParams := dbcm.getConnParams(key)
	return connParams.Name
}

func (dbcm *DBConnManager) connect(key string) string {

}
