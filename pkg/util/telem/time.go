package telem

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

/// |||| STAMP |||

func NewTimeStamp(time time.Time) TimeStamp {
	return TimeStamp(time.UnixMicro())
}

type TimeStamp int64

func (ts TimeStamp) Add(s TimeSpan) TimeStamp {
	return TimeStamp(int64(ts) + int64(s))
}

func (ts TimeStamp) String() string {
	return time.UnixMicro(int64(ts)).String()
}

const SecondsToMicroSeconds = 1000000

// |||| SPAN ||||

type TimeSpan int64

func NewTimeSpan(duration time.Duration) TimeSpan {
	return TimeSpan(duration.Microseconds())
}

func (ts TimeSpan) ToDataRate() DataRate {
	return DataRate(float64(SecondsToMicroSeconds) / float64(ts))
}

func (ts TimeSpan) ToDuration() time.Duration {
	return time.Duration(ts) * time.Microsecond
}

func (ts TimeSpan) String() string {
	return ts.ToDuration().String()
}

// |||| DATA RATE |||

type DataRate float64

func (dr DataRate) Period() TimeSpan {
	return TimeSpan(1 / float64(dr) * SecondsToMicroSeconds)
}

func (dr DataRate) String() string {
	fStr := strconv.FormatFloat(float64(dr), 'E', -1, 32)
	return fmt.Sprintf("%s Hz", fStr)
}

// |||| TIME RANGE ||||

type TimeRange struct {
	start TimeStamp
	end   TimeStamp
}

func NewTimeRange(start TimeStamp, end TimeStamp) TimeRange {
	return TimeRange{start: start, end: end}
}

func (tr TimeRange) Start() TimeStamp {
	return tr.start
}

func (tr TimeRange) End() TimeStamp {
	return tr.end
}

func (tr TimeRange) Span() TimeSpan {
	return TimeSpan(tr.end - tr.start)
}

func (tr TimeRange) Overlap(cTr TimeRange) (TimeRange, bool) {
	olStart := TimeStamp(math.Max(float64(tr.Start()), float64(cTr.Start())))
	olEnd := TimeStamp(math.Min(float64(tr.End()), float64(cTr.End())))
	olR := NewTimeRange(olStart, olEnd)
	return olR, olStart < olEnd
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("from %s to %s", tr.Start(), tr.End())
}
