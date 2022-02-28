package telem

import (
	"encoding/binary"
	"github.com/arya-analytics/aryacore/pkg/models"
	"math"
)

type Chunk struct {
	start TimeStamp
	dr    DataRate
	dt    models.ChannelDataType
	data  *ChunkData
}

func (t *Chunk) Start() TimeStamp {
	return t.start
}

func (t *Chunk) End() TimeStamp {
	return t.Start().Add(t.Span())

}

func (t *Chunk) Range() TimeRange {
	return NewTimeRange(t.Start(), t.End())
}

func (t *Chunk) Span() TimeSpan {
	return TimeSpan(t.Len() * int64(t.Period()))
}

func (t *Chunk) Len() int64 {
	return t.data.Size() / t.SampleSize()

}

func (t *Chunk) RemoveFromStart(ts TimeStamp) {
	t.data.Splice(t.ByteIndexAtTS(t.Start()), t.ByteIndexAtTS(ts))
	t.start = ts
}

func (t *Chunk) RemoveFromEnd(ts TimeStamp) {
	t.data.Splice(t.ByteIndexAtTS(ts), t.ByteIndexAtTS(t.End()))

}

func (t *Chunk) Duration() {
}

func (t *Chunk) SampleSize() int64 {
	switch t.dt {
	case models.ChannelDataTypeFloat64:
		return 8
	default:
		panic("t chunk has unknown data type")
	}
}

func (t *Chunk) Period() TimeSpan {
	return TimeSpan(1 / float64(t.dr) * SecondsToMicroSeconds)
}

func (t *Chunk) RangeSinceStart(ts TimeStamp) TimeRange {
	return NewTimeRange(t.start, ts)
}

func (t *Chunk) IndexAtTS(ts TimeStamp) int64 {
	return int64(t.RangeSinceStart(ts).Span() / t.Period())
}

func (t *Chunk) ByteIndexAtTS(ts TimeStamp) int64 {
	return t.IndexAtTS(ts) * t.SampleSize()
}

func (t *Chunk) ValueAtTS(ts TimeStamp) interface{} {
	byteI := t.ByteIndexAtTS(ts)
	return convertBytesToValue(t.data.ReadSlice(byteI, byteI+t.SampleSize()), t.dt)
}

func convertBytesToValue(b []byte, dataType models.ChannelDataType) interface{} {
	switch dataType {
	case models.ChannelDataTypeFloat64:
		bits := binary.BigEndian.Uint64(b)
		return math.Float64frombits(bits)
	default:
		panic("telem chunk has unknown data type")
	}
}
