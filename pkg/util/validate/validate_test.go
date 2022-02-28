package validate_test

import (
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	Describe("With No opts", func() {
		It("Should return the correct error", func() {
			v := validate.New[string]([]func(string) error{
				func(a string) error {
					return errors.New("a strange error")
				},
			})
			err := v.Exec("string").Error()
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(errors.New("a strange error")))
			Expect(v.Errors()).To(Equal([]error{errors.New("a strange error")}))
		})
	})
	Describe("With aggregation", func() {
		It("Should return all the errors", func() {
			v := validate.New[string]([]func(string) error{
				func(a string) error {
					return errors.New("a strange error")
				},
				func(a string) error {
					return errors.New("an even stranger error")
				},
			}, validate.WithAggregation())
			err := v.Exec("string").Error()
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(errors.New("a strange error")))
			Expect(v.Errors()).To(Equal([]error{errors.New("a strange error"), errors.New("an even stranger error")}))
		})
	})
})
