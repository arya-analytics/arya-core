package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx     = context.Background()
	driver  = mock.NewDriverRoach(false, false)
	pool    = storage.NewPool()
	engine  = roach.New(driver, pool)
	adapter = engine.NewAdapter()
)

var _ = BeforeSuite(func() {
	migrateErr := engine.NewMigrate().Exec(ctx)
	Expect(migrateErr).To(BeNil())
})

var _ = AfterSuite(func() {
	driver.Stop()
})

func TestRoach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roach Suite")
}
