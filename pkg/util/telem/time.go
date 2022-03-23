package telem

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

/// |||| STAMP |||

// TimeStamp is used to store timing values for telemetry. Internally, it represents an (unvalidated) int64
// representing a UTC Unix timestamp in microseconds.
type TimeStamp int64

const (
	// TimeStampMin represents the minimum possible timestamp.
	TimeStampMin TimeStamp = math.MinInt64
	// TimeStampMax represents the maximum possible time stamp.
	TimeStampMax TimeStamp = math.MaxInt64
)

// NewTimeStamp creates a new TimeStamp from a time.Time.
func NewTimeStamp(time time.Time) TimeStamp {
	return TimeStamp(time.UnixMicro())
}

// Add adds TimeSpan to a TimeStamp.
func (ts TimeStamp) Add(s TimeSpan) TimeStamp {
	return TimeStamp(int64(ts) + int64(s))
}

// ToTime converts tTimeStamp to a time.Time.
func (ts TimeStamp) ToTime() time.Time {
	return time.UnixMicro(int64(ts))
}

func (ts TimeStamp) String() string {
	return time.UnixMicro(int64(ts)).String()
}

// SecondsToMicroseconds converts a time/duration in seconds to microseconds.
const SecondsToMicroseconds = 1000000

// |||| SPAN ||||

// TimeSpan represents a span of time in Unix microseconds.
type TimeSpan int64

// NewTimeSpan creates a new TimeSpan from a time.Duration.
func NewTimeSpan(duration time.Duration) TimeSpan {
	return TimeSpan(duration.Microseconds())
}

// ToDataRate converts TimeSpan into a DataRate.
//
// Ex: a TimeSpan of 100 microseconds would be a data rate of 10 KHz
func (ts TimeSpan) ToDataRate() DataRate {
	return DataRate(float64(SecondsToMicroseconds) / float64(ts))
}

// ToDuration converts TimeSpan to a time.Duration.
func (ts TimeSpan) ToDuration() time.Duration {
	return time.Duration(ts) * time.Microsecond
}

func (ts TimeSpan) String() string {
	return ts.ToDuration().String()
}

// |||| DATA RATE |||

// DataRate stores the sample rate of telemetry (in Hz).
type DataRate float64

// Period returns a TimeSpan representing the amount of time between samples.
//
// Ex: a DataRate of 10 KHz would convert to a period of 100 microseconds.
//
func (dr DataRate) Period() TimeSpan {
	return TimeSpan(1 / float64(dr) * SecondsToMicroseconds)
}

func (dr DataRate) String() string {
	var precision int = 0
	if dr < 1 {
		precision = 3
	}
	fStr := strconv.FormatFloat(float64(dr), 'f', precision, 32)
	return fmt.Sprintf("%sHz", fStr)
}

// |||| TIME RANGE ||||

// TimeRange represents a range of time between two TimeStamp.
type TimeRange struct {
	start TimeStamp
	end   TimeStamp
}

// NewTimeRange constructs a TimeRange from a start and end TimeStamp.
func NewTimeRange(start TimeStamp, end TimeStamp) TimeRange {
	return TimeRange{start: start, end: end}
}

// AllTime returns the widest possible time range. TimeStampMin to TimeStampMax.
func AllTime() TimeRange {
	return NewTimeRange(TimeStampMin, TimeStampMax)
}

// Start returns a TimeStamp representing the start of TimeRange.
func (tr TimeRange) Start() TimeStamp {
	return tr.start
}

// End returns a TimeRange representing the end of TimeRange.
func (tr TimeRange) End() TimeStamp {
	return tr.end
}

// Span returns a TimeSpan representing the amount of time within TimeRange.
func (tr TimeRange) Span() TimeSpan {
	return TimeSpan(tr.end - tr.start)
}

// IsZero returns true if the range has a span of zero.
func (tr TimeRange) IsZero() bool {
	return tr.Span() == 0
}

// Overlap calculates and returns an overlap TimeRange between two TimeRange. Returns a second argument that is true
// if the overlap is valid, and false if the overlap isn't (i.e. the two TimeRange don't overlap).
func (tr TimeRange) Overlap(cTr TimeRange) (TimeRange, bool) {
	olStart := TimeStamp(math.Max(float64(tr.Start()), float64(cTr.Start())))
	olEnd := TimeStamp(math.Min(float64(tr.End()), float64(cTr.End())))
	olR := NewTimeRange(olStart, olEnd)
	return olR, olStart < olEnd
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("from %s to %s", tr.Start(), tr.End())
}
