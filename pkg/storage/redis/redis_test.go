package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redis Engine", func() {
	Describe("adapter", func() {
		var a pool.Adapt[*redis.Engine]
		BeforeEach(func() {
			var err error
			a, err = engine.NewAdapt(nil)
			Expect(err).To(BeNil())
		})
		Describe("New adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(a.Healthy()).To(BeTrue())
			})
		})
	})
})
