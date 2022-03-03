package storage_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	ctx   = context.Background()
	store *mock.Storage
)

var _ = BeforeSuite(func() {
	store = mock.NewStorage()
	Expect(store.NewMigrate().Exec(ctx)).To(BeNil())
})

var _ = AfterSuite(func() {
	store.Stop()
})

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "storage Suite")
}
