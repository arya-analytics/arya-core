package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	storagemock "github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx   = context.Background()
	clust *mock.Cluster
)

var _ = BeforeSuite(func() {
	var err error
	clust, err = mock.New(ctx,
		storagemock.WithVerbose())
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	clust.Stop()
})

func TestChanchunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chanchunk Suite")
}
