package cluster_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var cl = mock.NewCluster()

func TestCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster Suite")
}
