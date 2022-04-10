package streamq_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Stream", func() {
	Describe("Complete and Wait", func() {
		It("Should block until the query is complete", func() {
			streamQ := &streamq.Stream{Done: make(chan struct{})}
			go func() {
				time.Sleep(50 * time.Millisecond)
				streamQ.Complete()
			}()
			t0 := time.Now()
			streamQ.Wait()
			Expect(time.Since(t0)).To(BeNumerically(">=", 50*time.Millisecond))
		})
	})
	Describe("Segment", func() {
		It("Should start a new goroutine", func() {
			streamQ := &streamq.Stream{Done: make(chan struct{}), Segments: map[streamq.Segment]bool{}}
			i := 0
			streamQ.Segment(func() {
				time.Sleep(50 * time.Millisecond)
				i++
				streamQ.Complete()
			}, streamq.WithSegmentName("mysegment"))
			streamQ.Wait()
			Expect(i).To(Equal(1))
		})
	})

})
