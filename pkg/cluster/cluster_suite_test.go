package cluster_test

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
	Expect(store.NewMigrate().Exec(ctx)).To(BeNil())
})

var _ = AfterSuite(func() {
	store.Stop()
})

func TestCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster Suite")
}
