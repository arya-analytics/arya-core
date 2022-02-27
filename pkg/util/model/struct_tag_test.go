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
			It("Should retrieve the correct tag by kev:Val pair", func() {
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
		Describe("HasAnyFields", func() {
			DescribeTable("Should return the correct boolean", func(expected bool, flds ...string) {
				Expect(tags.HasAnyFields(flds...)).To(Equal(expected))
			},
				Entry("false for an empty field", false, ""),
				Entry("false for a field that doesn't exist on the model", false, "RandomField"),
				Entry("false for a nested field that doesn't exist on the model", false, "Field.Field"),
				Entry("true for a nested field whose first field exists on the model", true, "ID.RandomField"),
				Entry("true for a nested field whose field exists on the model", true, "InnerModel.ID"),
			)
		})
	})
	Describe("Edge cases + errors", func() {
		Context("No category provided", func() {
			It("Should return false", func() {
				_, ok := tags.Retrieve("*", "random", "hello")
				Expect(ok).To(BeFalse())
			})
		})
		Context("Val but no key provided", func() {
			It("Should panic", func() {
				Expect(func() {
					tags.Retrieve("randomcat", "*", "Val")
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
