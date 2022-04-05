package query_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Opt", func() {
	Describe("OptConvertChain", func() {
		It("Should execute the convert chain correctly", func() {
			v := 0
			occ := query.OptConvertChain{func(p *query.Pack) {
				v = 1
			}}
			Expect(func() {
				occ.Exec(query.NewRetrieve().Pack())
			}).ToNot(Panic())
			Expect(v).To(Equal(1))
		})
	})
})
