package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockDB     testserver.TestServer
	mockEngine *roach.Engine
	mockNode   = &storage.Node{
		ID: 1,
	}
	mockModel = &storage.ChannelConfig{
		ID:     uuid.New(),
		Name:   "Cool Name",
		NodeID: mockNode.ID,
	}
	mockCtx     = context.Background()
	mockADapter *roach.Adapter
)

func migrate() {
	err := mockEngine.NewMigrate(mockADapter).Verify(mockCtx)
	if err != nil {
		if err := mockEngine.NewMigrate(mockADapter).Exec(mockCtx); err != nil {
			log.Fatalln(err)
		}
	}
}

func createMockModel() {
	if err := mockEngine.NewCreate(mockADapter).Model(mockModel).Exec(mockCtx); err != nil {
		log.Fatalln(err)
	}
}

func deleteMockModel() {
	if err := mockEngine.NewDelete(mockADapter).Model(mockModel).WherePK(
		mockModel.ID).Exec(mockCtx); err != nil {
		log.Fatalln(err)
	}
}

func createMockNode() {
	if err := mockEngine.NewCreate(mockADapter).Model(mockNode).Exec(
		mockCtx); err != nil {
		log.Fatalln(err)
	}
}

func bootstrapEngine() {
	var err error
	mockDB, err = testserver.NewTestServer()
	if err != nil {
		log.Fatalln(err)
	}
	mockEngine = roach.New(roach.Config{DSN: mockDB.PGURL().String(), Driver: roach.DriverPG})
	mockADapter = mockEngine.NewAdapter().(*roach.Adapter)
}

var _ = BeforeSuite(func() {
	bootstrapEngine()
	migrate()
	createMockNode()
})

var _ = AfterSuite(func() {
	mockDB.Stop()
})

func TestRoach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roach Suite")
}
