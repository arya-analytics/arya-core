package caseconv_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Caseconv", func() {
	Context("PascalToKebab", func() {
		It("Should correctly convert the string", func() {
			Expect(caseconv.PascalToKebab("HelloStrange")).To(Equal("hello-strange"))
		})
	})
	Context("PascalToSnake", func() {
		It("Shouild correctly convert the string", func() {
			Expect(caseconv.PascalToSnake("HelloStrange")).To(Equal("hello_strange"))
		})
	})
})
