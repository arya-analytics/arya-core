package mock

import (
	"database/sql"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"net"
	"net/url"
	"strconv"
)

// |||| ROACH ||||

type DriverRoach struct {
	Host     string
	Port     int
	HTTPPort int
	Username string
	Password string
}

func availHTTPPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return listener.Addr().(*net.TCPAddr).Port
}

func (d *DriverRoach) Connect() (*bun.DB, error) {
	port := availHTTPPort()
	d.HTTPPort = port
	ts, err := testserver.NewTestServer(testserver.HTTPPortOpt(port),
		testserver.SecureOpt(), testserver.RootPasswordOpt("testpass"))
	if err != nil {
		return nil, err
	}
	sqlDB, err := sql.Open("postgres", ts.PGURL().String())
	if cErr := d.bindConnProperties(ts.PGURL()); cErr != nil {
		return nil, cErr
	}
	if err != nil {
		return nil, err
	}
	return bun.NewDB(sqlDB, pgdialect.New()), nil
}

func (d *DriverRoach) bindConnProperties(url *url.URL) error {
	d.Host = url.Hostname()
	port, err := strconv.Atoi(url.Port())
	if err != nil {
		return err
	}
	d.Port = port
	uname := url.User.Username()
	d.Username = uname
	pwd, ok := url.User.Password()
	if !ok {
		return errors.New("could not bind password")
	}
	d.Password = pwd
	return nil
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
