package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roach", func() {
	var a pool.Adapt[*roach.Engine]
	BeforeEach(func() {
		var err error
		a, err = engine.NewAdapt(nil)
		Expect(err).To(BeNil())
	})
	Describe("Adapter", func() {
		Describe("New adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(a.Healthy()).To(BeTrue())
			})
		})
	})
})
