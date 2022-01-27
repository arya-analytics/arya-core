package timeseries_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockBaseClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	mockClient = timeseries.NewWrap(mockBaseClient)
	mockCtx    = context.Background()
	mockTSKey  = "mockTSKey"
)

func createMockTS() {
	if err := mockClient.TSCreate(mockCtx, mockTSKey, timeseries.CreateOptions{
		Retention: 0,
	}).Err(); err != nil {
		log.Fatalln(err)
	}
}

func deleteMockTS() {
	if err := mockClient.Del(mockCtx, mockTSKey).Err(); err != nil {
		log.Fatalln(err)
	}
}

func TestTimeseries(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timeseries Suite")
}
