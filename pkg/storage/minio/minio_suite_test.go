package minio_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockEngine  = minio.New(mock.DriverMinio{})
	mockAdapter = mockEngine.NewAdapter()
	mockCtx     = context.Background()
	mockBytes   = []byte("mock model bytes")
	mockModel   *storage.ChannelChunk
)

func createMockModel() {
	mockModel = &storage.ChannelChunk{
		ID:   uuid.New(),
		Data: mock.NewObject([]byte("mock model bytes")),
	}
	if err := mockEngine.NewCreate(mockAdapter).Model(mockModel).Exec(
		mockCtx); err != nil {
		log.Fatalln(err)
	}
}

func deleteMockModel() {
	if err := mockEngine.NewDelete(mockAdapter).Model(mockModel).WherePK(mockModel.ID).
		Exec(
			mockCtx); err != nil {
		log.Fatalln(err)
	}
}

func TestMinio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
