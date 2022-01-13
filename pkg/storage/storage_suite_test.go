package storage_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"testing"
)


var (
	dummyEngineCfg = storage.EngineConfig{
		storage.EngineRoleMetaData: &roach.Engine{
			Driver: roach.DriverSQLite,
		},
	}
	dummyStorage = storage.New(dummyEngineCfg)
	dummyCtx = context.Background()
	dummyModel = &storage.ChannelConfig{
		ID:   432,
		Name: "Cool Name",
	}
)

func createDummyModel() {
	if err := dummyStorage.NewCreate().Model(dummyModel).Exec(
		dummyCtx); err != nil {
		log.Fatalln(err)
	}
}

func deleteDummyModel() {
	if err := dummyStorage.NewDelete().Model(dummyModel).WhereID(dummyModel.ID).Exec(
		dummyCtx); err != nil {
		log.Fatalln(err)
	}
}


func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()
	if err := dummyStorage.Migrate(ctx); err != nil {
		log.Fatalln(err)
	}
})

