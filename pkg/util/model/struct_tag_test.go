package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("StructTag", func() {
	tags := model.NewReflect(&mock.ModelA{}).StructTagChain()
	Describe("Standard Usage", func() {
		Describe("Retrieving a tag", func() {
			It("Should retrieve the correct tag by category", func() {
				t, ok := tags.Retrieve("randomcat", "*", "*")
				Expect(ok).To(BeTrue())
				Expect(t.Field.Name).To(Equal("Name"))
			})
			It("Should retrieve the correct tag by key", func() {
				t, ok := tags.Retrieve("model", "role", "*")
				Expect(ok).To(BeTrue())
				Expect(t.Field.Name).To(Equal("ID"))
			})
			It("Should retrieve the correct tag by kev:value pair", func() {
				t, ok := tags.Retrieve("randomcat", "random", "hello")
				Expect(ok).To(BeTrue())
				Expect(t.Field.Name).To(Equal("Name"))
			})
			It("Should return false if the tag doesnt exist", func() {
				_, ok := tags.Retrieve("randomcat", "random", "lalala")
				Expect(ok).To(BeFalse())
			})
			It("Should retrieve the correct tag by role", func() {
				_, ok := tags.RetrieveByFieldRole(model.PKRole)
				Expect(ok).To(BeTrue())
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("No category provided", func() {
			It("Should return false", func() {
				_, ok := tags.Retrieve("*", "random", "hello")
				Expect(ok).To(BeFalse())
			})
		})
		Context("value but no key provided", func() {
			It("Should panic", func() {
				Expect(func() {
					tags.Retrieve("randomcat", "*", "value")
				}).To(Panic())
			})
		})
		Context("Non-struct provided to constructor", func() {
			It("Should panic", func() {
				Expect(func() {
					model.NewStructTagChain(reflect.TypeOf(1))
				}).To(Panic())
			})
		})
	})
})
