package streamq_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TSCreate", func() {
	var (
		exec = &mock.Exec{}
		asm  = streamq.NewAssembleTS(exec.Exec)
	)
	Describe("BindStream", func() {
		It("Should bind the stream correctly", func() {
			streamQ := &streamq.Stream{
				Errors: make(chan error),
				Ctx:    nil,
			}
			p := asm.NewTSCreate().BindStream(streamQ).Pack()
			so, ok := streamq.RetrieveStreamOpt(p)
			Expect(ok).To(BeTrue())
			Expect(so).To(Equal(streamQ))
		})
	})
	Describe("Stream", func() {
		It("Should return the stream", func() {
			ctx := context.Background()
			stream, err := asm.NewTSCreate().Stream(ctx)
			Expect(err).To(BeNil())
			Expect(stream.Ctx).To(Equal(ctx))
		})
	})
})
