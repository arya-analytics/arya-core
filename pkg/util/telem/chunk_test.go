package telem_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ChunkData", func() {
	Describe("The Basics", func() {

		var (
			dataSize = 10000
			data     []float64
			startTs  = telem.TimeStamp(time.Now().UnixMicro())
			chunk    *telem.Chunk
		)
		BeforeEach(func() {
			data = []float64{}
			for i := 0; i < dataSize; i++ {
				data = append(data, float64(i))
			}

			bulk := telem.NewChunkData([]byte{})
			Expect(bulk.WriteData(data)).To(BeNil())
			chunk = telem.NewChunk(startTs, telem.DataTypeFloat64, telem.DataRate(25), bulk)
		})
		Describe("Static Attributes", func() {
			Describe("Sample Size", func() {
				It("Should return the correct sample size", func() {
					Expect(chunk.SampleSize()).To(Equal(int64(8)))
				})
			})
			Describe("Period", func() {
				It("Should return the correct period", func() {
					Expect(chunk.Period()).To(Equal(telem.TimeSpan(40000)))
				})
			})

		})
		Describe("Timing", func() {
			Describe("RangeFromStart", func() {
				It("Should return the correct start", func() {
					Expect(chunk.Start()).To(Equal(startTs))
				})
				It("Should return the correct end", func() {
					Expect(chunk.End()).To(Equal(startTs.Add(telem.NewTimeSpan(time.Duration(dataSize/25) * time.Second))))
				})
				It("Should return the correct range", func() {
					ts := chunk.Start().Add(chunk.Period() * 30)
					Expect(chunk.RangeFromStart(ts).Span()).To(Equal(chunk.Period() * 30))
					Expect(chunk.RangeFromStart(ts).Start()).To(Equal(chunk.Start()))
					Expect(chunk.RangeFromStart(ts).End()).To(Equal(ts))
				})
			})
		})
		Describe("Len", func() {
			It("Should return the correct length", func() {
				Expect(chunk.Len()).To(Equal(int64(len(data))))
			})
		})
		Describe("Span", func() {
			It("Should return the correct span", func() {
				Expect(chunk.Span()).To(Equal(telem.TimeSpan(400000000)))
			})
		})
		Describe("Range", func() {
			It("Should return the correct range", func() {
				Expect(chunk.Range().Start()).To(Equal(chunk.Start()))
				Expect(chunk.Range().End()).To(Equal(chunk.End()))
				Expect(chunk.Range().Span()).To(Equal(chunk.Span()))
			})
		})
		Describe("IndexAt", func() {
			It("Should return the correct sample index", func() {
				ts := chunk.Start().Add(chunk.Period() * 30)
				Expect(chunk.IndexAt(ts)).To(Equal(int64(30)))
			})
			It("Should return the correct sample byte index", func() {
				ts := chunk.Start().Add(chunk.Period() * 30)
				Expect(chunk.ByteIndexAt(ts)).To(Equal(int64(30 * 8)))
			})
		})
		Describe("ValueAT", func() {
			It("Should return the correct value at the timestamp", func() {
				ts := chunk.Start().Add(chunk.Period() * 30)
				Expect(chunk.ValueAt(ts).(float64)).To(Equal(float64(30)))
			})
		})
		Describe("Remove", func() {
			It("Should removeFrom from the start", func() {
				ts := chunk.Start().Add(chunk.Period() * 30)
				chunk.RemoveFromStart(ts)
				Expect(chunk.Size())
				Expect(chunk.ValueAt(chunk.Start())).To(Equal(float64(30)))
			})
			It("Should removeFrom from the end", func() {
				ts := chunk.Start().Add(chunk.Period() * 400)
				chunk.RemoveFromEnd(ts)
				Expect(chunk.Span()).To(Equal(chunk.Period() * 400))
				Expect(chunk.ValueAt(chunk.Start())).To(Equal(float64(0)))
				Expect(chunk.ValueAt(chunk.End())).To(Equal(float64(400)))
			})
		})
	})
	Describe("Time Correctness Checks", func() {
		var (
			c *telem.Chunk
		)
		BeforeEach(func() {
			cd := telem.NewChunkData([]byte{})
			Expect(cd.WriteData([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})).To(BeNil())
			c = telem.NewChunk(
				telem.TimeStamp(0),
				telem.DataTypeFloat64,
				telem.DataRate(1),
				cd,
			)
		})
		Describe("Times", func() {
			It("Should return the correct start time", func() {
				Expect(c.Start()).To(Equal(telem.TimeStamp(0)))
			})
			It("Should return the correct end time", func() {
				Expect(c.End()).To(Equal(c.Start().Add(telem.NewTimeSpan(9 * time.Second))))
			})
			It("Should return the correct span", func() {
				Expect(c.Span()).To(Equal(telem.NewTimeSpan(9 * time.Second)))
			})
		})
		Describe("Value Access", func() {
			Describe("Single Value", func() {
				It("Should return the correct value at the first timestamp", func() {
					Expect(c.ValueAt(telem.TimeStamp(0))).To(Equal(1.0))
				})
				It("Should panic when trying to access the value at the last timestamp", func() {
					Expect(func() {
						c.ValueAt(c.End())
					}).To(Panic())
				})
				It("Should return the correct last value", func() {
					Expect(c.ValueAt(c.End().Add(telem.NewTimeSpan(-1 * time.Second)))).To(Equal(9.0))
				})
			})
			Describe("Range of SourceValues", func() {
				It("Should return the correct values", func() {
					rng := c.RangeFromStart(c.Start().Add(telem.NewTimeSpan(3 * time.Second)))
					Expect(rng.Start()).To(Equal(c.Start()))
					Expect(rng.End()).To(Equal(c.Start().Add(telem.NewTimeSpan(3 * time.Second))))
					vals := c.ValuesInRange(rng).([]float64)
					Expect(vals).To(Equal([]float64{1, 2, 3}))
				})
				Describe("Exceeding the data range", func() {
					It("Should return an empty slice when there is no overlap", func() {
						start := c.Start().Add(telem.NewTimeSpan(11 * time.Second))
						end := c.Start().Add(telem.NewTimeSpan(20 * time.Second))
						rng := telem.NewTimeRange(start, end)
						vals := c.ValuesInRange(rng).([]float64)
						Expect(vals).To(Equal([]float64{}))
					})
				})
				Describe("Consuming the data range", func() {
					It("Should return all values in the range", func() {
						start := c.Start().Add(telem.NewTimeSpan(-20 * time.Second))
						end := c.Start().Add(telem.NewTimeSpan(20 * time.Second))
						rng := telem.NewTimeRange(start, end)
						vals := c.ValuesInRange(rng).([]float64)
						Expect(vals).To(Equal([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
					})
				})
			})
			Describe("All Values", func() {
				It("Should return all values", func() {
					Expect(c.AllValues()).To(Equal([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
				})
			})
			Describe("Large Chunk", func() {
				It("Shouldn't take an absurd amount of time to convert the values", func() {
					cd := telem.NewChunkData([]byte{})
					var vals []float64
					for i := 0; i < 20000; i++ {
						vals = append(vals, float64(i))
					}
					Expect(cd.WriteData(vals)).To(BeNil())
					cc := telem.NewChunk(telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(25000), cd)
					ts := time.Now()
					cc.AllValues()
					Expect(time.Since(ts)).To(BeNumerically("<", 5*time.Millisecond))
				})
			})
		})
	})

})
