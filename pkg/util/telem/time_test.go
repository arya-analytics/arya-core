package telem_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
	"time"
)

var _ = Describe("Time", func() {
	Describe("TimeStamp", func() {
		It("Should create a new timestamp from a time", func() {
			t := time.Now()
			ts := telem.NewTimeStamp(t)
			Expect(int64(ts)).To(Equal(t.UnixMicro()))
		})
		It("Should correctly convert a timestamp to a time", func() {
			t := time.UnixMilli(10)
			ts := telem.NewTimeStamp(t)
			Expect(ts.ToTime()).To(Equal(time.UnixMilli(10)))
		})
		It("Should stringify the timestamp", func() {
			t := time.UnixMilli(10)
			ts := telem.NewTimeStamp(t)
			Expect(reflect.TypeOf(ts.String())).To(Equal(reflect.TypeOf("")))
		})
	})
	Describe("TimeSpan", func() {
		It("Should create a new span from a duration", func() {
			d := 1 * time.Second
			ts := telem.NewTimeSpan(d)
			Expect(ts.ToDuration()).To(Equal(d))
		})
		It("Should return the span as a data rate", func() {
			d := 1 * time.Second
			ts := telem.NewTimeSpan(d)
			Expect(ts.ToDataRate()).To(Equal(telem.DataRate(1)))
		})
		It("Should stringify the timespan", func() {
			ts := telem.NewTimeSpan(10 * time.Second)
			Expect(ts.String()).To(Equal("10s"))
		})
	})
	Describe("TimeRange", func() {
		It("Should create the correct time range", func() {
			t0 := time.Now()
			t1 := time.Now().Add(1 * time.Second)
			rng := telem.NewTimeRange(telem.NewTimeStamp(t0), telem.NewTimeStamp(t1))
			Expect(int64(rng.Start())).To(Equal(t0.UnixMicro()))
			Expect(int64(rng.End())).To(Equal(t1.UnixMicro()))
			Expect(rng.Span()).To(BeNumerically(">", telem.NewTimeSpan(900*time.Millisecond)))
		})
		It("Should stringify the time range correctly", func() {
			t0 := time.UnixMilli(10)
			t1 := time.UnixMilli(20)
			rng := telem.NewTimeRange(telem.NewTimeStamp(t0), telem.NewTimeStamp(t1))
			Expect(len(rng.String())).To(Equal(len("from 1969-12-31 19:00:00.01 -0500 EST to 1969-12-31 19:00:00.02 -0500 EST")))
		})
		Describe("ChunkOverlap", func() {
			Context("No overlap", func() {
				Describe("Clearly no overlap", func() {
					It("Should return false", func() {
						baseT := time.Now()
						rngOne := telem.NewTimeRange(
							telem.NewTimeStamp(baseT.Add(2*time.Second)),
							telem.NewTimeStamp(baseT.Add(3*time.Second)),
						)
						rngTwo := telem.NewTimeRange(
							telem.NewTimeStamp(baseT),
							telem.NewTimeStamp(baseT.Add(1*time.Second)),
						)
						_, oOneOverlap := rngOne.Overlap(rngTwo)
						Expect(oOneOverlap).To(BeFalse())
						_, oTwoOverlap := rngTwo.Overlap(rngOne)
						Expect(oTwoOverlap).To(BeFalse())
					})
				})
				Describe("Send and end times are the same", func() {
					It("Should return no overlap", func() {
						baseT := time.Now()
						rngOne := telem.NewTimeRange(
							telem.NewTimeStamp(baseT.Add(1*time.Second)),
							telem.NewTimeStamp(baseT.Add(2*time.Second)),
						)
						rngTwo := telem.NewTimeRange(
							telem.NewTimeStamp(baseT),
							telem.NewTimeStamp(baseT.Add(1*time.Second)),
						)
						_, oOneOverlap := rngOne.Overlap(rngTwo)
						Expect(oOneOverlap).To(BeFalse())
						_, oTwoOverlap := rngTwo.Overlap(rngOne)
						Expect(oTwoOverlap).To(BeFalse())
					})
				})
			})
			Context("Partial ChunkOverlap", func() {
				It("Should return the correct overlap range", func() {
					baseT := time.Now()
					rngOne := telem.NewTimeRange(
						telem.NewTimeStamp(baseT.Add(1*time.Second)),
						telem.NewTimeStamp(baseT.Add(4*time.Second)),
					)
					rngTwo := telem.NewTimeRange(
						telem.NewTimeStamp(baseT),
						telem.NewTimeStamp(baseT.Add(3*time.Second)),
					)
					oOneOverlap, oOneOverlapExists := rngOne.Overlap(rngTwo)
					Expect(oOneOverlapExists).To(BeTrue())
					Expect(oOneOverlap.Span()).To(Equal(telem.NewTimeSpan(2 * time.Second)))
					Expect(oOneOverlap.End()).To(Equal(rngTwo.End()))
					Expect(oOneOverlap.Start()).To(Equal(rngOne.Start()))

					oTwoOverlap, oTwoOverlapExists := rngTwo.Overlap(rngOne)
					Expect(oTwoOverlapExists).To(BeTrue())
					Expect(oTwoOverlap.Span()).To(Equal(telem.NewTimeSpan(2 * time.Second)))
					Expect(oTwoOverlap.End()).To(Equal(rngTwo.End()))
					Expect(oTwoOverlap.Start()).To(Equal(rngOne.Start()))
				})
			})
			Context("One TimeStampExp Inside Another", func() {
				It("Should return the correct overlap range", func() {
					baseT := time.Now()
					rngOne := telem.NewTimeRange(
						telem.NewTimeStamp(baseT),
						telem.NewTimeStamp(baseT.Add(4*time.Second)),
					)
					rngTwo := telem.NewTimeRange(
						telem.NewTimeStamp(baseT.Add(1*time.Second)),
						telem.NewTimeStamp(baseT.Add(3*time.Second)),
					)
					oOneOverlap, oOneOverlapExists := rngOne.Overlap(rngTwo)
					Expect(oOneOverlapExists).To(BeTrue())
					Expect(oOneOverlap.Span()).To(Equal(telem.NewTimeSpan(2 * time.Second)))
					Expect(oOneOverlap.End()).To(Equal(rngTwo.End()))
					Expect(oOneOverlap.Start()).To(Equal(rngTwo.Start()))

					oTwoOverlap, oTwoOverlapExists := rngOne.Overlap(rngTwo)
					Expect(oTwoOverlapExists).To(BeTrue())
					Expect(oTwoOverlap.Span()).To(Equal(telem.NewTimeSpan(2 * time.Second)))
					Expect(oTwoOverlap.End()).To(Equal(rngTwo.End()))
					Expect(oTwoOverlap.Start()).To(Equal(rngTwo.Start()))
				})
			})
		})
	})
	Describe("DataRate", func() {
		It("Should stringify the data rate correctly", func() {
			Expect(telem.DataRate(25).String()).To(Equal("25Hz"))
			Expect(telem.DataRate(0.5).String()).To(Equal("0.500Hz"))
		})
	})
})
