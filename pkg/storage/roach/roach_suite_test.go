package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx     = context.Background()
	engine  = roach.New(mock.DriverPG{})
	adapter = engine.NewAdapter()
)

var _ = BeforeSuite(func() {
	migrateErr := engine.NewMigrate(adapter).Exec(ctx)
	Expect(migrateErr).To(BeNil())
})

func TestRoach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roach Suite")
}
