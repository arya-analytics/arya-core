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
	store = mock.NewStorage()
)

var _ = BeforeSuite(func() {
	err := store.NewMigrate().Exec(ctx)
	Expect(err).To(BeNil())
})

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}
