package tsquery_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Retrieve", func() {
	var (
		exec = &mock.Exec{}
		asm  = tsquery.NewAssemble(exec.Exec)
	)
	Describe("TimeRangeOpt", func() {
		It("Should set the time range all opt properly", func() {
			p := asm.NewRetrieve().AllTime().Pack()
			tr, ok := tsquery.TimeRangeOpt(p)
			Expect(ok).To(BeTrue())
			Expect(tr.Start()).To(Equal(telem.TimeRangeMin))
			Expect(tr.End()).To(Equal(telem.TimeRangeMax))
		})
		It("Should set the where time range opt properly", func() {
			tr := telem.NewTimeRange(telem.NewTimeStamp(time.Now()), telem.NewTimeStamp(time.Now()))
			p := asm.NewRetrieve().WhereTimeRange(tr).Pack()
			opt, ok := tsquery.TimeRangeOpt(p)
			Expect(ok).To(BeTrue())
			Expect(opt).To(Equal(tr))
		})
		It("Should return ok as false when the opt isn't specified", func() {
			p := asm.NewRetrieve().Pack()
			_, ok := tsquery.TimeRangeOpt(p)
			Expect(ok).To(BeFalse())
		})

	})
})
