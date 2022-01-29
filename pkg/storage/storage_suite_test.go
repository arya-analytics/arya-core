package storage_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"testing"
)

var (
	mockEngineCfg = storage.EngineConfig{
		storage.EngineRoleMD: bootstrapMockRoachEngine(),
		storage.EngineRoleObject: minio.New(
			minio.Config{
				Driver:    minio.DriverMinIO,
				Endpoint:  "play.min.io",
				AccessKey: "Q3AM3UQ867SPQQA43P2F",
				SecretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			}),
		storage.EngineRoleCache: redis.New(redis.Config{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			Database: 0,
		}),
	}
	mockStorage    = storage.New(mockEngineCfg)
	mockCtx        = context.Background()
	mockChannelCfg = &storage.ChannelConfig{
		ID:     uuid.New(),
		Name:   "Cool Name",
		NodeID: 1,
	}
	mockBytes = []byte("mock model bytes")
	mockNode  = &storage.Node{
		ID: 1,
	}
	mockRange = &storage.Range{
		ID:                uuid.New(),
		LeaseHolderNodeID: mockNode.ID,
	}
	mockChannelChunk *storage.ChannelChunk
)

func bootstrapMockRoachEngine() storage.MDEngine {
	var err error
	mockDB, err := testserver.NewTestServer()
	if err != nil {
		log.Fatalln(err)
	}
	return roach.New(roach.Config{DSN: mockDB.PGURL().String(),
		Driver: roach.DriverPG})
}

func createMock(m interface{}) {
	rfl := model.NewReflect(m)
	if err := mockStorage.NewCreate().Model(m).Exec(mockCtx); err != nil {
		log.Fatalln(err, rfl.Type().Name())
	}
}

func deleteMock(m interface{}) {
	rfl := model.NewReflect(m)
	if err := mockStorage.NewDelete().Model(m).WherePK(rfl.PKField().
		Value().Interface()).Exec(mockCtx); err != nil {
		log.Fatalln(err, rfl.Type().Name())
	}
}

func createMockChannelCfg() {
	createMock(mockChannelCfg)
}

func deleteMockChannelCfg() {
	deleteMock(mockChannelCfg)
}

func createMockRange() {
	createMock(mockRange)
}

func deleteMockRange() {
	deleteMock(mockRange)
}

func createMockChannelChunk() {
	createMockChannelCfg()
	createMockRange()
	mockChannelChunk = &storage.ChannelChunk{
		ID:              uuid.New(),
		Data:            mock.NewObject(mockBytes),
		RangeID:         mockRange.ID,
		ChannelConfigID: mockChannelCfg.ID,
	}
	createMock(mockChannelChunk)
}

func deleteMockChannelChunk() {
	deleteMockChannelCfg()
	deleteMockRange()
	deleteMock(mockChannelChunk)
}

func createMockNode() {
	createMock(mockNode)
}

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()
	if err := mockStorage.NewMigrate().Exec(ctx); err != nil {
		log.Fatalln(err)
	}
	createMockNode()
})
