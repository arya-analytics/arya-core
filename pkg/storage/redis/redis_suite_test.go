package redis_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	ctx    = context.Background()
	engine *redis.Engine
)

var _ = BeforeSuite(func() {
	engine = redis.New(mock.DriverRedis{})
})

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Suite")
}
