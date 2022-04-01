package redis_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	ctx    = context.Background()
	engine internal.EngineCache
)

var _ = BeforeSuite(func() {
	p := pool.New[internal.Engine]()
	engine = redis.New(mock.DriverRedis{}, p)
	p.AddFactory(engine)
})

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Suite")
}
