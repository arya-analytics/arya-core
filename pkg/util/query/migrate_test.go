package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrate", func() {
	var (
		exec = &mock.Exec{}
		asm  = query.NewAssemble(exec.Exec)
	)
	Describe("Exec", func() {
		It("Should execute the migrations correctly", func() {
			Expect(asm.NewMigrate().Exec(ctx)).To(BeNil())
		})
	})
	Describe("Verify", func() {
		It("Should add the option correctly", func() {
			p := asm.NewMigrate().Verify().Pack()
			v := query.RetrieveVerifyOpt(p)
			Expect(v).To(BeTrue())
		})
		It("Should be false when the option isn't specified", func() {
			p := asm.NewMigrate().Pack()
			v := query.RetrieveVerifyOpt(p)
			Expect(v).To(BeFalse())
		})
	})
})
