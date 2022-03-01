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
	var asm = query.NewAssemble(&mock.Exec{})
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
})
