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
	"github.com/uptrace/bun/extra/bundebug"
	"io/ioutil"
	baseLog "log"
	"math"
	"net"
	"net/url"
	"strconv"
	"time"
)

// |||| ROACH ||||

type DriverRoach struct {
	Host     string
	Port     int
	WithHTTP bool
	HTTPPort int
	Username string
	Password string
	Verbose  bool
	servers  []testserver.TestServer
}

func NewDriverRoach(withHTTP bool, verbose bool) *DriverRoach {
	return &DriverRoach{WithHTTP: withHTTP, Verbose: verbose}
}

func availHTTPPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return listener.Addr().(*net.TCPAddr).Port
}

func (d *DriverRoach) Connect() (*bun.DB, error) {
	baseLog.SetOutput(ioutil.Discard)
	if d.WithHTTP {
		d.HTTPPort = availHTTPPort()
	}
	ts, err := testserver.NewTestServer(
		testserver.HTTPPortOpt(d.HTTPPort),
		testserver.SecureOpt(),
		testserver.RootPasswordOpt("testpass"),
	)

	wErr := ts.WaitForInit()
	if wErr != nil {
		return nil, wErr
	}
	d.servers = append(d.servers, ts)
	if err != nil {
		return nil, err
	}
	sqlDB, sErr := sql.Open("postgres", ts.PGURL().String())
	if sErr != nil {
		return nil, sErr
	}
	if cErr := d.bindConnProperties(ts.PGURL()); cErr != nil {
		return nil, cErr
	}
	if err != nil {
		return nil, err
	}
	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	if d.Verbose {
		bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return bunDB, nil
}

func (d *DriverRoach) DemandCap() int {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to return math.MaxInt.
	return math.MaxInt
}

func (d *DriverRoach) Expiration() time.Duration {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to make it huuuge.
	return 1000000 * time.Hour
}

func (d *DriverRoach) Healthy() bool {
	// Connections never fail!
	return true
}

func (d *DriverRoach) Stop() {
	for _, ts := range d.servers {
		ts.Stop()
	}
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
	pwd, _ := url.User.Password()
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

func (d DriverRedis) DemandCap() int {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to return math.MaxInt.
	return math.MaxInt
}

func (d DriverRedis) Expiration() time.Duration {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to make it huuuge.
	return 1000000 * time.Hour
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

func (d DriverMinio) DemandCap() int {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to return math.MaxInt.
	return math.MaxInt
}

func (d DriverMinio) Expiration() time.Duration {
	// Because Connect creates a new database, we never want to
	// recycle the connection, so we need to make it huuuge.
	return 1000000 * time.Hour
}
