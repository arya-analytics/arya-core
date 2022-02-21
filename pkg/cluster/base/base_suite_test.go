package base_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx   = context.Background()
	store *mock.Storage
)

var _ = BeforeSuite(func() {
	store = mock.NewStorage()
	err := store.NewMigrate().Exec(ctx)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	store.Stop()
})

func TestBase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}