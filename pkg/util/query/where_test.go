package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Where", func() {
	DescribeTable("Field Expressions", func(exp query.FieldExp, expOp query.FieldOp, expVals []interface{}) {
		Expect(exp.Op).To(Equal(expOp))
		Expect(exp.Vals).To(HaveLen(len(expVals)))
		Expect(exp.Vals).To(Equal(expVals))
	},
		Entry("Greater Than", query.GreaterThan(1), query.FieldOpGreaterThan, []interface{}{1}),
		Entry("Less Than", query.LessThan(1), query.FieldOpLessThan, []interface{}{1}),
		Entry("In Range", query.InRange(1, 2), query.FieldOpInRange, []interface{}{1, 2}),
		Entry("In", query.In(1, 2, 3), query.FieldOpIn, []interface{}{1, 2, 3}),
	)
	Describe("Opt Binding", func() {
		var asm query.Assemble
		BeforeEach(func() {
			asm = query.NewAssemble(&mock.Exec{})
		})
		Describe("Primary Keys", func() {
			It("Should create the correct pk opt", func() {
				p := asm.NewRetrieve().WherePK(123).Pack()
				pkc, ok := query.PKOpt(p)
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(1))
				Expect(pkc[0].Raw()).To(Equal(123))
			})
			It("Should create the correct multi pk opt", func() {
				p := asm.NewRetrieve().WherePKs([]int{1, 2, 3}).Pack()
				pkc, ok := query.PKOpt(p)
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(3))
				Expect(pkc[0].Raw()).To(Equal(1))
			})
			It("Should return false when a pk opt wasn't specified", func() {
				p := asm.NewRetrieve().Pack()
				pkc, ok := query.PKOpt(p)
				Expect(pkc).To(HaveLen(0))
				Expect(ok).To(BeFalse())
			})
			DescribeTable("Single primary keys on Query Variants", func(q query.Query) {
				pkc, ok := query.PKOpt(q.Pack())
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(1))
			},
				Entry("Delete", asm.NewDelete().WherePK(1)),
				Entry("Update", asm.NewUpdate().WherePK(1)),
			)
			DescribeTable("Multi primary keys on Query Variants", func(q query.Query) {
				pkc, ok := query.PKOpt(q.Pack())
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(3))
			},
				Entry("Delete", asm.NewDelete().WherePKs([]int{1, 2, 3})),
			)

		})
		Describe("Where Fields", func() {
			It("Should create the correct where fields opt", func() {
				p := asm.NewRetrieve().WhereFields(query.WhereFields{"RandomField": "RandomValue"}).Pack()
				wf, ok := query.WhereFieldsOpt(p)
				Expect(ok).To(BeTrue())
				Expect(wf["RandomField"]).To(Equal("RandomValue"))
				Expect(len(wf)).To(Equal(1))
			})
			It("Should return false when a where fields opt wasn't specified", func() {
				p := asm.NewRetrieve().Pack()
				wf, ok := query.WhereFieldsOpt(p)
				Expect(ok).To(BeFalse())
				Expect(wf).To(HaveLen(0))
			})
		})
	})
})
