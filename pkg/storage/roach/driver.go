package roach

import (
	"database/sql"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	// DSN is a connection string for the database. If specified,
	// all other fields except for Driver can be left blank.
	DSN string
	// Username for the database.
	Username string
	// Password for the database.
	Password string
	// Host IP for the database.
	Host string
	// Port to Connect to at Host.
	Port int
	// Database to Connect to.
	Database string
	// Whether to open a TLS connection or not.
	UseTLS bool
	// TransactionLogLevel
	TransactionLogLevel TransactionLogLevel
}

const (
	driverRoachDSN                 = "dsn"
	driverRoachUsername            = "username"
	driverRoachPassword            = "password"
	driverRoachHost                = "host"
	driverRoachPort                = "port"
	driverRoachDatabase            = "database"
	driverRoachUseTLS              = "useTLS"
	driverRoachTransactionLogLevel = "transactionLogLevel"
)

func configLeaf(name string) string {
	return storage.ConfigTree().MetaData().Leaf(name)
}

func (c Config) Viper() Config {
	return Config{
		DSN:                 viper.GetString(configLeaf(driverRoachDSN)),
		Username:            viper.GetString(configLeaf(driverRoachUsername)),
		Password:            viper.GetString(configLeaf(driverRoachPassword)),
		Host:                viper.GetString(configLeaf(driverRoachHost)),
		Port:                viper.GetInt(configLeaf(driverRoachPort)),
		Database:            viper.GetString(configLeaf(driverRoachDatabase)),
		UseTLS:              viper.GetBool(configLeaf(driverRoachUseTLS)),
		TransactionLogLevel: TransactionLogLevel(viper.GetInt(configLeaf(driverRoachTransactionLogLevel))),
	}
}

type DriverRoach struct {
	Config Config
}

func (d DriverRoach) Connect() (*bun.DB, error) {
	c := d.buildConnector()
	db := sql.OpenDB(c)
	bunDB := bun.NewDB(db, pgdialect.New())
	setLogLevel(d.Config.TransactionLogLevel, bunDB)
	return bunDB, nil
}

func (d DriverRoach) addr() string {
	return fmt.Sprintf("%s:%v", d.Config.Host, d.Config.Port)
}

func (d DriverRoach) buildConnector() *pgdriver.Connector {
	if d.Config.DSN != "" {
		return pgdriver.NewConnector(pgdriver.WithDSN(d.Config.DSN))
	}
	return pgdriver.NewConnector(
		pgdriver.WithAddr(d.addr()),
		pgdriver.WithInsecure(d.Config.UseTLS),
		pgdriver.WithUser(d.Config.Username),
		pgdriver.WithPassword(d.Config.Password),
		pgdriver.WithDatabase(d.Config.Database))
}

type TransactionLogLevel int

const (
	// TransactionLogLevelNone logs no queries.
	TransactionLogLevelNone TransactionLogLevel = iota
	// TransactionLogLevelErr logs failed queries.
	TransactionLogLevelErr
	// TransactionLogLevelAll logs all queries.
	TransactionLogLevelAll
)

func setLogLevel(t TransactionLogLevel, db *bun.DB) {
	switch t {
	case TransactionLogLevelAll:
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	case TransactionLogLevelErr:
		db.AddQueryHook(bundebug.NewQueryHook())
	}
}
