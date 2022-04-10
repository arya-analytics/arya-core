package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec", func() {
	It("Should do nothing", func() {
		e := &mock.Exec{}
		Expect(e.Exec(context.Background(), query.NewRetrieve().Pack())).To(Succeed())
	})

})
