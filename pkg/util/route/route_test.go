package route_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	Describe("CtxDone", func() {
		It("Should return true if the context is cancelled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			Expect(route.CtxDone(ctx)).To(BeTrue())
		})
		It("Should return false if the context is not cancelled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			Expect(route.CtxDone(ctx)).To(BeFalse())
			cancel()
		})
	})
	Describe("RangeContext", func() {
		It("Should range through the values in the channel until the context is cancelled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			var (
				resV []int
				inv  = make(chan int)
			)
			go func() {
				defer cancel()
				for i := 0; i < 5; i++ {
					inv <- i
				}
			}()
			route.RangeContext(ctx, inv, func(v int) {
				resV = append(resV, v)
			})
			Expect(resV).To(Equal([]int{0, 1, 2, 3, 4}))
		})
		It("Should range through the values until the channel is closed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			var (
				resV []int
				inv  = make(chan int)
			)
			go func() {
				defer close(inv)
				for i := 0; i < 5; i++ {
					inv <- i
				}
			}()
			route.RangeContext(ctx, inv, func(v int) {
				resV = append(resV, v)
			})
			Expect(resV).To(Equal([]int{0, 1, 2, 3, 4}))
			cancel()
		})
	})
})
