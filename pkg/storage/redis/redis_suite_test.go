package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	mockEngine := redis.New(redis.Config{
		Host:     "redis-13401.c278.us-east-1-4.ec2.cloud.redislabs.com",
		Port:     13401,
		Password: "6GNc75OA8q2IUAjymRqq16ZemK2OoBQ9",
	})
	_ = mockEngine.NewAdapter()
})

func TestRedists(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redists Suite")
}
