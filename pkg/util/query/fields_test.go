package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fields", func() {
	var (
		asm query.Assemble
	)
	BeforeEach(func() {
		asm = query.NewAssemble(&mock.Exec{})
	})
	Describe("FieldsOpt", func() {
		var (
			fo query.FieldsOpt
		)
		BeforeEach(func() {
			p := asm.NewRetrieve().Fields("One", "Two", "Three").Pack()
			var ok bool
			fo, ok = query.RetrieveFieldsOpt(p)
			Expect(ok).To(BeTrue())
		})
		It("Should return false if the opt isn't specified", func() {
			p := asm.NewRetrieve().Pack()
			foT, ok := query.RetrieveFieldsOpt(p)
			Expect(ok).To(BeFalse())
			Expect(foT).To(HaveLen(0))
		})
		Describe("ContainsAny", func() {
			It("Should return false when the opt doesn't contain the fields", func() {
				Expect(fo.ContainsAny("Four")).To(BeFalse())
			})
			It("Should return true when the opt contains the fields", func() {
				Expect(fo.ContainsAny("Three")).To(BeTrue())
			})
		})
		Describe("AllExcept", func() {
			It("Should return all the fields expect for the field specified", func() {
				Expect(fo.AllExcept("One").ContainsAny("One")).To(BeFalse())
			})
		})
		Describe("ContainsAll", func() {
			It("Should return false when one of the fields isn't in the opt", func() {
				Expect(fo.ContainsAll("One", "Seven")).To(BeFalse())
			})
			It("Should return ture when all of the fields are in the opt", func() {
				Expect(fo.ContainsAll("One", "Three")).To(BeTrue())
			})
		})
		Describe("Append", func() {
			It("Should append a result to the opt", func() {
				Expect(fo.Append("Four").ContainsAll("Four")).To(BeTrue())
			})
			It("Shouldn't append any duplicates", func() {
				Expect(fo.Append("Three")).To(HaveLen(3))
			})
		})
	})
})
