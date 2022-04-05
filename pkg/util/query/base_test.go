package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	modelmock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Base", func() {
	var (
		exec = &mock.Exec{}
		asm  = query.NewAssemble(exec.Exec)
	)
	Describe("Common Query Functionality", func() {
		DescribeTable("Model",
			func(q query.Query) {
				p := q.Pack()
				Expect(p.Model().Type()).To(Equal(model.NewReflect(&modelmock.ModelA{}).Type()))
			},
			Entry("Create", asm.NewCreate().Model(&modelmock.ModelA{})),
			Entry("Retrieve", asm.NewRetrieve().Model(&modelmock.ModelA{})),
			Entry("Update", asm.NewUpdate().Model(&modelmock.ModelA{})),
			Entry("Delete", asm.NewDelete().Model(&modelmock.ModelA{})),
		)
	})
	Describe("Exec", func() {
		It("Should execute the query", func() {
			Expect(asm.Exec(ctx, query.NewMigrate().Pack())).To(BeNil())
		})
		It("Should panic if no execute is bound", func() {
			Expect(func() { query.NewRetrieve().Exec(ctx) }).To(Panic())
		})
	})
})
