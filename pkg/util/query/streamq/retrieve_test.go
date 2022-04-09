package streamq_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("TSRetrieve", func() {
	var (
		exec = &mock.Exec{}
		asm  = streamq.NewAssembleTS(exec.Exec)
	)
	Describe("RetrieveTimeRangeOpt", func() {
		It("Should set the time range all opt properly", func() {
			p := asm.NewTSRetrieve().AllTime().Pack()
			tr, ok := streamq.RetrieveTimeRangeOpt(p)
			Expect(ok).To(BeTrue())
			Expect(tr.Start()).To(Equal(telem.TimeStampMin))
			Expect(tr.End()).To(Equal(telem.TimeStampMax))
		})
		It("Should set the where time range opt properly", func() {
			tr := telem.NewTimeRange(telem.NewTimeStamp(time.Now()), telem.NewTimeStamp(time.Now()))
			p := asm.NewTSRetrieve().WhereTimeRange(tr).Pack()
			opt, ok := streamq.RetrieveTimeRangeOpt(p)
			Expect(ok).To(BeTrue())
			Expect(opt).To(Equal(tr))
		})
		It("Should return ok as false when the opt isn't specified", func() {
			p := asm.NewTSRetrieve().Pack()
			_, ok := streamq.RetrieveTimeRangeOpt(p)
			Expect(ok).To(BeFalse())
		})

	})
})
