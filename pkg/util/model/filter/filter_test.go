package filter_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model/filter"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("By Primary Key", func() {
		It("Should filter by pk correctly", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}}
			var os []*mock.ModelA
			filter.Exec(query.NewRetrieve().Model(&s).WherePK(1).Pack(), &os)
			Expect(s).To(HaveLen(3))
			Expect(os).To(HaveLen(1))
		})
		It("Should maintain duplicates", func() {
			s := []*mock.ModelA{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 1}}
			var os []*mock.ModelA
			filter.Exec(query.NewRetrieve().Model(&s).WherePK(1).Pack(), &os)
			Expect(s).To(HaveLen(4))
			Expect(os).To(HaveLen(2))
		})
	})
})
