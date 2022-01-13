package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyEngine = &roach.Engine{
	Driver: roach.DriverSQLite,
}

func TestRoach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roach Suite")
}
