package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Minio Engine", func() {
	Describe("adapter", func() {
		var a pool.Adapt[internal.Engine]
		BeforeEach(func() {
			var err error
			a, err = engine.NewAdapt()
			Expect(err).To(BeNil())
		})
		Describe("New adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(a.Healthy()).To(BeTrue())
			})
		})
	})
})
