package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	dummyEngine = roach.New(roach.Config{
		//Host:                "192.168.64.11",
		//Port:                26257,
		//Username:            "root",
		//Password:            "",
		//Database:            "defaultdb",
		TransactionLogLevel: roach.TransactionLogLevelAll,
		Driver:              roach.DriverSQLite,
	})
	dummyNode = &storage.Node{
		ID: 1,
	}
	dummyModel = &storage.ChannelConfig{
		ID:     uuid.New(),
		Name:   "Cool Name",
		NodeID: dummyNode.ID,
	}
	dummyCtx     = context.Background()
	dummyAdapter = dummyEngine.NewAdapter()
)

func migrate() {
	err := dummyEngine.NewMigrate(dummyAdapter).Verify(dummyCtx)
	if err != nil {
		if err := dummyEngine.NewMigrate(dummyAdapter).Exec(dummyCtx); err != nil {
			log.Fatalln(err)
		}
	}
}

func createDummyModel() {
	if err := dummyEngine.NewCreate(dummyAdapter).Model(dummyModel).Exec(dummyCtx); err != nil {
		log.Fatalln(err)
	}
}

func deleteDummyModel() {
	if err := dummyEngine.NewDelete(dummyAdapter).Model(dummyModel).WherePK(
		dummyModel.ID).Exec(dummyCtx); err != nil {
		log.Fatalln(err)
	}
}

var _ = BeforeSuite(migrate)

func TestRoach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roach Suite")
}
