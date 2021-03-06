package clusterapi_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mockDriver *mock.DriverRoach

var _ = BeforeSuite(func() {
	mockDriver = mock.NewDriverRoach(true, false)
	_, err := mockDriver.Connect()
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	mockDriver.Stop()
})

func TestClusterAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClusterAPI Suite")
}
