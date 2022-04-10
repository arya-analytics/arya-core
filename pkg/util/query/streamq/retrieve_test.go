package streamq_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
	Describe("WherePK", func() {
		It("Should set the primary key properly", func() {
			p := asm.NewTSRetrieve().WherePK(1).Pack()
			pkc, ok := query.RetrievePKOpt(p)
			Expect(ok).To(BeTrue())
			Expect(pkc).To(HaveLen(1))
			Expect(pkc[0].Raw()).To(Equal(1))
		})
	})
	Describe("WherePKs", func() {
		It("Should set the primary keys correctly", func() {
			p := asm.NewTSRetrieve().WherePKs([]int{1, 2, 3}).Pack()
			pkc, ok := query.RetrievePKOpt(p)
			Expect(ok).To(BeTrue())
			Expect(pkc).To(HaveLen(3))
			Expect(pkc[0].Raw()).To(Equal(1))
			Expect(pkc[1].Raw()).To(Equal(2))
			Expect(pkc[2].Raw()).To(Equal(3))
		})
	})
	Describe("BindStream", func() {
		It("Should bind the stream correctly", func() {
			streamQ := &streamq.Stream{
				Errors: make(chan error),
				Ctx:    nil,
			}
			p := asm.NewTSRetrieve().BindStream(streamQ).Pack()
			so, ok := streamq.RetrieveStreamOpt(p)
			Expect(ok).To(BeTrue())
			Expect(so).To(Equal(streamQ))
		})
	})
	Describe("Stream", func() {
		It("Should return the stream", func() {
			ctx := context.Background()
			stream, err := asm.NewTSRetrieve().Stream(ctx)
			Expect(err).To(BeNil())
			Expect(stream.Ctx).To(Equal(ctx))
		})
	})
})
