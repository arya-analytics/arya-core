package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", func() {
	var (
		asm  query.Assemble
		exec *mock.Exec
	)
	BeforeEach(func() {
		exec = &mock.Exec{}
		asm = query.NewAssemble(exec.Exec)
	})
	Describe("Bulk", func() {
		It("Should set the bulk update opt correctly", func() {
			Expect(asm.NewUpdate().Bulk().Exec(ctx)).To(BeNil())
			bulk := query.BulkUpdateOpt(exec.Pack)
			Expect(bulk).To(BeTrue())
		})
		It("Should return false when the opt isn't specified", func() {
			Expect(asm.NewUpdate().Exec(ctx)).To(BeNil())
			bulk := query.BulkUpdateOpt(exec.Pack)
			Expect(bulk).To(BeFalse())
		})
		It("Should panic when trying to retrieve a bulk opt from a non update", func() {
			Expect(asm.NewRetrieve().Exec(ctx)).To(BeNil())
			Expect(func() {
				query.BulkUpdateOpt(exec.Pack)
			}).To(Panic())
		})
	})
	Describe("Fields", func() {
		It("Should set the correct fields", func() {
			p := asm.NewUpdate().Fields("One", "Two").Pack()
			fo, ok := query.RetrieveFieldsOpt(p)
			Expect(ok).To(BeTrue())
			Expect(fo).To(HaveLen(2))
		})
	})
})
