package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Where", func() {
	DescribeTable("Field Expressions", func(exp query.FieldExp, expOp query.FieldOp, expVals []interface{}) {
		Expect(exp.Op).To(Equal(expOp))
		Expect(exp.Values).To(HaveLen(len(expVals)))
		Expect(exp.Values).To(Equal(expVals))
	},
		Entry("Greater Than", query.GreaterThan(1), query.FieldOpGreaterThan, []interface{}{1}),
		Entry("Less Than", query.LessThan(1), query.FieldOpLessThan, []interface{}{1}),
		Entry("In Range", query.InRange(1, 2), query.FieldOpInRange, []interface{}{1, 2}),
		Entry("In", query.In(1, 2, 3), query.FieldOpIn, []interface{}{1, 2, 3}),
	)
	Describe("Opt Binding", func() {
		var (
			exec = &mock.Exec{}
			asm  = query.NewAssemble(exec.Exec)
		)
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
				Expect(pkc).To(Equal(model.NewPKChain([]int{1, 2, 3})))
			})
			It("Should return false when a pk opt wasn't specified", func() {
				p := asm.NewRetrieve().Pack()
				pkc, ok := query.PKOpt(p)
				Expect(pkc).To(HaveLen(0))
				Expect(ok).To(BeFalse())
			})
			It("Should allow the caller to pass in a PKChain", func() {
				inPKC := model.NewPKChain([]uuid.UUID{uuid.New(), uuid.New()})
				p := asm.NewRetrieve().WherePKs(inPKC).Pack()
				pkc, ok := query.PKOpt(p)
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(2))
				Expect(pkc[0]).To(Equal(inPKC[0]))
				Expect(pkc).To(Equal(inPKC))
			})
			It("Should allow the caller to pass in a PK", func() {
				inPK := model.NewPK(uuid.New())
				p := asm.NewRetrieve().WherePK(inPK).Pack()
				pkc, ok := query.PKOpt(p)
				Expect(ok).To(BeTrue())
				Expect(pkc).To(HaveLen(1))
				Expect(pkc[0].Raw()).To(Equal(inPK.Raw()))
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
			Describe("Invalid primary keys", func() {
				It("Should panic when the caller passes a slice of priamry keys to WherePK", func() {
					Expect(func() {
						asm.NewRetrieve().WherePK([]int{1, 2, 3})
					}).To(Panic())
				})
			})

		})
		Describe("Where Fields", func() {
			It("Should create the correct Where fields opt", func() {
				p := asm.NewRetrieve().WhereFields(query.WhereFields{"RandomField": "RandomValue"}).Pack()
				wf, ok := query.RetrieveWhereFieldsOpt(p)
				Expect(ok).To(BeTrue())
				Expect(wf["RandomField"]).To(Equal("RandomValue"))
				Expect(len(wf)).To(Equal(1))
			})
			It("Should return false when a Where fields opt wasn't specified", func() {
				p := asm.NewRetrieve().Pack()
				wf, ok := query.RetrieveWhereFieldsOpt(p)
				Expect(ok).To(BeFalse())
				Expect(wf).To(HaveLen(0))
			})
		})
	})
})
