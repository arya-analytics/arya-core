package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx    = context.Background()
	driver = mock.NewDriverRoach(false, false)
	p      = pool.New[internal.Engine]()
	engine = roach.New(driver, p)
)

var _ = BeforeSuite(func() {
	p.AddFactory(engine)
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
