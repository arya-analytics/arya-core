package tsquery_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	modelmock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Assemble", func() {
	var (
		exec = &mock.Exec{}
		asm  = tsquery.NewAssemble(exec.Exec)
	)
	Describe("Common Query Functionality", func() {
		DescribeTable("Model",
			func(q query.Query) {
				p := q.Pack()
				Expect(p.Model().Type()).To(Equal(model.NewReflect(&modelmock.ModelA{}).Type()))
			},
			Entry("Create", asm.NewCreate().Model(&modelmock.ModelA{})),
			Entry("Retrieve", asm.NewRetrieve().Model(&modelmock.ModelA{})),
		)
	})
	Describe("Exec", func() {
		It("Should execute the query correctly", func() {
			Expect(asm.Exec(context.Background(), tsquery.NewRetrieve().Pack()))
		})
	})
})
