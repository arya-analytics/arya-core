package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retrieve", func() {
	var (
		exec = &mock.Exec{}
		asm  = query.NewAssemble(exec.Exec)
	)
	Describe("Calc", func() {
		var (
			into  *int
			field = "field"
		)
		DescribeTable(
			"Different Calc",
			func(op query.Calc) {
				p := asm.NewRetrieve().Calc(op, field, into).Pack()
				calcOpt, ok := query.RetrieveCalcOpt(p)
				Expect(ok).To(BeTrue())
				Expect(calcOpt.Op).To(Equal(op))
				Expect(calcOpt.Field).To(Equal(field))
				Expect(calcOpt.Into).To(Equal(into))
			},
			Entry("Sum", query.CalcSum),
			Entry("Max", query.CalcMax),
			Entry("Min", query.CalcMin),
			Entry("Count", query.CalcCount),
			Entry("Avg", query.CalcAVG),
		)
		Describe("It should return false when the calc opt isn't specified", func() {
			p := asm.NewRetrieve().Pack()
			_, ok := query.RetrieveCalcOpt(p)
			Expect(ok).To(BeFalse())
		})
		It("Should panic when into is not a pointer", func() {
			Expect(func() { asm.NewRetrieve().Calc(query.CalcSum, "Field", 1) }).To(Panic())
		})
	})
	Describe("Relation", func() {
		It("Should create a relation opt correctly", func() {
			p := asm.NewRetrieve().Relation("Name", "Fld").Pack()
			ro := query.RetrieveRelationOpts(p)
			Expect(ro).To(HaveLen(1))
			Expect(ro[0].Name).To(Equal("Name"))
			Expect(ro[0].Fields).To(Equal(query.FieldsOpt{"Fld"}))
		})
		It("Should allow for multiple relation opts", func() {
			p := asm.NewRetrieve().
				Relation("Name", "Fld").
				Relation("RelTwo", "FldTwo").
				Pack()
			ro := query.RetrieveRelationOpts(p)
			Expect(ro).To(HaveLen(2))
		})
		It("Should return an empty array when no relations are specified", func() {
			p := asm.NewRetrieve().Pack()
			ro := query.RetrieveRelationOpts(p)
			Expect(ro).To(HaveLen(0))
		})
	})
	Describe("Fields", func() {
		It("Should create a fields opt correctly", func() {
			p := asm.NewRetrieve().Fields("FieldOne", "FieldTwo").Pack()
			fo, ok := query.RetrieveFieldsOpt(p)
			Expect(ok).To(BeTrue())
			Expect(fo).To(HaveLen(2))
			Expect(fo).To(Equal(query.FieldsOpt{"FieldOne", "FieldTwo"}))
		})
	})
	Describe("Limit", func() {
		It("Should create the limit opt correctly", func() {
			p := asm.NewRetrieve().Limit(1).Pack()
			limit, ok := query.RetrieveLimitOpt(p)
			Expect(ok).To(BeTrue())
			Expect(limit).To(Equal(1))
		})
		It("Should return false if no limit is specified", func() {
			p := asm.NewRetrieve().Pack()
			limit, ok := query.RetrieveLimitOpt(p)
			Expect(ok).To(BeFalse())
			Expect(limit).To(Equal(0))
		})
	})
	Describe("OrderDirection", func() {
		It("Should set an order opt on the correct field", func() {
			p := asm.NewRetrieve().Order(query.OrderASC, "Field").Pack()
			order, ok := query.RetrieveOrderOpt(p)
			Expect(ok).To(BeTrue())
			Expect(order.Field).To(Equal("Field"))
			Expect(order.Direction).To(Equal(query.OrderASC))
		})
		It("Should return false if no order is specified", func() {
			p := asm.NewRetrieve().Pack()
			_, ok := query.RetrieveOrderOpt(p)
			Expect(ok).To(BeFalse())
		})

	})
})
