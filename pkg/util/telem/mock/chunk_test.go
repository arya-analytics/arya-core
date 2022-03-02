package mock_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Chunk", func() {
	Describe("PopulateRand", func() {
		It("Should populate the chunk with random float64 vals", func() {
			cd := telem.NewChunkData([]byte{})
			mock.PopulateRand(cd, telem.DataTypeFloat64, 100)
			c := telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(1), cd)
			Expect(c.Len()).To(Equal(int64(100)))
		})
		It("Should populate the chunk with random float32 vals", func() {
			cd := telem.NewChunkData([]byte{})
			mock.PopulateRand(cd, telem.DataTypeFloat32, 100)
			c := telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat32, telem.DataRate(1), cd)
			Expect(c.Len()).To(Equal(int64(100)))
		})
	})
	Describe("Contiguous Chunks", func() {
		It("Should create a set of contiguous chunks", func() {
			cc := mock.ChunkSet(
				5,
				telem.TimeStamp(0),
				telem.DataTypeFloat32,
				telem.DataRate(25),
				telem.NewTimeSpan(30*time.Second),
				telem.TimeSpan(0),
			)
			Expect(cc).To(HaveLen(5))
			Expect(cc[0].Start()).To(Equal(telem.TimeStamp(0)))
			Expect(cc[1].Start()).To(Equal(telem.TimeStamp(0).Add(telem.NewTimeSpan(30 * time.Second))))
			for i := 1; i < len(cc); i++ {
				Expect(cc[i-1].End()).To(Equal(cc[i].Start()))
			}
			Expect(cc[0].Len()).To(Equal(int64(25 * 30)))
			Expect(cc[0].Span()).To(Equal(telem.NewTimeSpan(30 * time.Second)))
		})
	})
})
