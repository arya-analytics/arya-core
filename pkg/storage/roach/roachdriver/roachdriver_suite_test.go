package roachdriver_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mockDriver *mock.DriverRoach

var _ = BeforeSuite(func() {
	mockDriver = &mock.DriverRoach{}
	_, err := mockDriver.Connect()
	Expect(err).To(BeNil())
})

func TestRoachdriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roachdriver Suite")
}
