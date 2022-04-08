package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete", func() {
	var (
		exec = &mock.Exec{}
		asm  = query.NewAssemble(exec.Exec)
	)
	Describe("WhereFields", func() {
		p := asm.NewDelete().WhereFields(query.WhereFields{"Hello": "World"}).Pack()
		wf, ok := query.RetrieveWhereFieldsOpt(p)
		Expect(ok).To(BeTrue())
		Expect(wf["Hello"]).To(Equal("World"))
	})
})
