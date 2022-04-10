package filter_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/filter"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("By Primary Key", func() {
		It("Should filter by pk correctly", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			os, err := filter.Filter(query.NewRetrieve().WherePK(1).Pack(), s)
			Expect(err).To(BeNil())
			Expect(s).To(HaveLen(3))
			Expect(os).To(HaveLen(1))
		})
		It("Should maintain duplicates", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 1}}
			os, err := filter.Filter(query.NewRetrieve().WherePK(1).Pack(), s)
			Expect(err).To(BeNil())
			Expect(s).To(HaveLen(4))
			Expect(os).To(HaveLen(2))
		})
	})
	Describe("By WhereFields", func() {
		It("Should filter by WhereFields", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			os, err := filter.Filter(query.NewRetrieve().WhereFields(query.WhereFields{"ID": 2}).Pack(), s)
			Expect(err).To(BeNil())
			Expect(s).To(HaveLen(3))
			Expect(os).To(HaveLen(1))
		})
		It("Should not panic on nonexistent fields", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			os, err := filter.Filter(query.NewRetrieve().WhereFields(query.WhereFields{"IDontExist": 4}).Pack(), s)
			Expect(err).To(BeNil())
			Expect(os).To(HaveLen(0))
		})
		It("Should panic when trying to use a field expression", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			Expect(func() {
				_, _ = filter.Filter(query.NewRetrieve().WhereFields(query.WhereFields{"ID": query.GreaterThan(2)}).Pack(), s)
			}).To(Panic())
		})
	})
	Describe("Calc", func() {
		It("Should calculate the value correctly", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			var res int64
			_, err := filter.Filter(query.NewRetrieve().Calc(query.CalcSum, "ID", &res).Pack(), s)
			Expect(err).To(BeNil())
			Expect(res).To(Equal(int64(6)))
		})
		It("Should calculate a floating point value correctly", func() {
			s := []*mock.ModelA{{IDFloat64: 1}, {IDFloat64: 2}, {IDFloat64: 3}}
			var res float64
			_, err := filter.Filter(query.NewRetrieve().Calc(query.CalcSum, "IDFloat64", &res).Pack(), s)
			Expect(err).To(BeNil())
			Expect(res).To(Equal(6.0))
		})
		It("Should panic when doing an operation on a non number", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			var res string
			Expect(func() {
				_, _ = filter.Filter(query.NewRetrieve().Calc(query.CalcSum, "Name", &res).Pack(), s)
			}).To(Panic())
		})
		It("Should panic when doing an unsupported operation", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			var res string
			Expect(func() {
				_, _ = filter.Filter(query.NewRetrieve().Calc(query.CalcMax, "ID", &res).Pack(), s)
			}).To(Panic())
		})
	})
	Describe("Options", func() {
		Describe("ErrorOnNotFound", func() {
			It("Should return an error if not found", func() {
				s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
				_, err := filter.Filter(query.NewRetrieve().WherePK(4).Pack(), s, filter.ErrorOnNotFound())
				Expect(err).ToNot(BeNil())
			})
			It("Should not return an error if found", func() {
				s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
				_, err := filter.Filter(query.NewRetrieve().WherePK(2).Pack(), s, filter.ErrorOnNotFound())
				Expect(err).To(BeNil())
			})
		})
		Describe("ErrorOnMultiple", func() {
			It("Should return an error if multiple found", func() {
				s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
				_, err := filter.Filter(query.NewRetrieve().Pack(), s, filter.ErrorOnMultiple())
				Expect(err).ToNot(BeNil())
			})
			It("Should not return an error if found", func() {
				s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
				_, err := filter.Filter(query.NewRetrieve().WherePK(2).Pack(), s, filter.ErrorOnMultiple())
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("ReflectFilter", func() {
		It("Should panic when the reflection doesn't contain a slice", func() {
			s := &mock.ModelA{ID: 1}
			Expect(func() {
				filter.ReflectFilter(query.NewRetrieve().Pack(), model.NewReflect(s))
			}).To(Panic())
		})
	})
})
