package errutil_test

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	Describe("Catcher", func() {
		Context("No error encountered", func() {
			var catcher = &errutil.Catcher{}
			It("Should continue to execute functions", func() {
				counter := 1
				for i := 0; i < 4; i++ {
					catcher.Exec(func() error {
						counter++
						return nil
					})
				}
				Expect(counter).To(Equal(5))
			})
			It("Should contain a nil error", func() {
				Expect(catcher.Error()).To(BeNil())
			})
		})
		Context("Error encountered", func() {
			var catcher = &errutil.Catcher{}
			It("Should stop execution", func() {
				counter := 1
				for i := 0; i < 4; i++ {
					catcher.Exec(func() error {
						if i == 2 {
							return fmt.Errorf("encountered unknown error")
						}
						counter++
						return nil
					})
				}
				Expect(counter).To(Equal(3))
			})
			It("Should contain a non-nil error", func() {
				Expect(catcher.Error()).ToNot(BeNil())
			})
		})
	})

})
