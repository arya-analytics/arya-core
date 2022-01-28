package redis_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/google/uuid"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockEngine = redis.New(redis.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
	})
	mockAdapter = mockEngine.NewAdapter()
	mockCtx     = context.Background()
	mockSeries  *storage.ChannelConfig
)

func createMockSeries() {
	mockSeries = &storage.ChannelConfig{
		Name: "SG_02",
		ID:   uuid.New(),
	}
	err := mockEngine.NewTSCreate(mockAdapter).Series().Model(mockSeries).Exec(mockCtx)
	if err != nil {
		panic(err)
	}
}

func TestRedists(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redists Suite")
}
