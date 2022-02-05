package redis_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/google/uuid"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockEngine = redis.New(redis.DriverRedis{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
	})
	mockAdapter = mockEngine.NewAdapter()
	mockCtx     = context.Background()
	mockSeries  *storage.ChannelConfig
	mockSamples []*storage.ChannelSample
)

func createMockSeries() {
	mockSeries = &storage.ChannelConfig{
		Name: "SG_02",
		ID:   uuid.New(),
	}
	if err := mockEngine.NewTSCreate(mockAdapter).Series().Model(mockSeries).Exec(
		mockCtx); err != nil {
		panic(err)
	}
}

func createMockSamples(qty int) {
	createMockSeries()
	mockSamples = []*storage.ChannelSample{}
	for i := 0; i < qty; i++ {
		// Sleeping to ensure we get unique timestamps
		time.Sleep(1 * time.Millisecond)
		mockSamples = append(mockSamples,
			&storage.ChannelSample{
				Timestamp:       time.Now().UnixNano(),
				Value:           123.2,
				ChannelConfigID: mockSeries.ID,
			},
		)
	}
	if err := mockEngine.NewTSCreate(mockAdapter).Sample().Model(&mockSamples).Exec(
		mockCtx); err != nil {
		panic(err)
	}
}

func TestRedists(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redists Suite")
}
