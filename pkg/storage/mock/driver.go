package mock

import (
	"database/sql"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"io/ioutil"
	baseLog "log"
	"net"
	"net/url"
	"strconv"
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
		log.Infof("Roach Connection String: %s", ts.PGURL().String())
	}
	return bunDB, nil
}

func (d *DriverRoach) DemandCap() int {
	return 5000000
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
	return 50
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
	return 50
}
