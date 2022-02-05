package mock

import (
	"database/sql"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// |||| PG ||||

type DriverPG struct{}

func (d DriverPG) Connect() (*bun.DB, error) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		return nil, err
	}
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(ts.PGURL().String())))
	db := bun.NewDB(sqlDB, pgdialect.New())
	return db, nil
}

// |||| REDIS ||||

type DriverRedis struct{}

func (d DriverRedis) Connect() (*timeseries.Client, error) {
	return timeseries.NewWrap(redis.NewClient(d.buildConfig())), nil

}

func (d DriverRedis) buildConfig() *redis.Options {
	return &redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	}
}

// |||| MINIO ||||

type DriverMinio struct{}

func (d DriverMinio) Connect() (*minio.Client, error) {
	return minio.New("localhost:9000", d.buildConfig())
}

func (d DriverMinio) buildConfig() *minio.Options {
	return &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	}
}
