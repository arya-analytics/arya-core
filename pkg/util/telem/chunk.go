// Package bulktelem holds utilities for working with telemetry, including storage, transport, traversal, and modification.
package telem

// |||| DATA TYPE ||||

type DataType int32

//go:generate stringer -type=DataType
const (
	DataTypeFloat64 DataType = iota
	DataTypeFloat32
)

// |||| CHUNK ||||

// Chunk wraps the binary telemetry container ChunkData and adds functionality for
// traversing it as if it was a slice of concrete values. Provides utilities for
// calculating ChunkOverlap, TimeRange, TimeSpan, values at indexes, and modifying
// the underlying ChunkData.
//
// Avoid instantiating directly, and instead call NewChunk.
type Chunk struct {
	start    TimeStamp
	DataType DataType
	DataRate DataRate
	*ChunkData
}

// NewChunk instantiates a new Chunk with the provided parameters.
func NewChunk(start TimeStamp, dataType DataType, dataRate DataRate, data *ChunkData) *Chunk {
	return &Chunk{start: start, DataType: dataType, DataRate: dataRate, ChunkData: data}
}

// |||| STATIC ATTRIBUTES ||||

// SampleSize return the size (in bytes) of a sample in the chunk.
func (t *Chunk) SampleSize() int64 {
	return sampleSize(t.DataType)
}

// Period returns a TimeSpan representing the period of time between samples.
func (t *Chunk) Period() TimeSpan {
	return t.DataRate.Period()
}

// |||| SIZING ||||

// Len returns the number of samples in the chunk.
func (t *Chunk) Len() int64 {
	return t.Size() / t.SampleSize()
}

// |||| TIMING ||||

// RangeFromStart returns a TimeRange representing the amount of time between the Start of the chunk and the
// provided TimeStamp ts.
func (t *Chunk) RangeFromStart(ts TimeStamp) TimeRange {
	return NewTimeRange(t.Start(), ts)
}

// Start returns a TimeStamp representing the start of the chunk.
func (t *Chunk) Start() TimeStamp {
	return t.start
}

// End returns a TimeStamp representing the end of the chunk.
func (t *Chunk) End() TimeStamp {
	return t.Start().Add(t.Span())

}

// Range returns the TimeRange between the Start and End of the chunk.
func (t *Chunk) Range() TimeRange {
	return NewTimeRange(t.Start(), t.End())
}

// Span returns a TimeSpan representing the time the data spans.
func (t *Chunk) Span() TimeSpan {
	return TimeSpan(t.Len() * int64(t.DataRate.Period()))
}

// |||| INDEXING ||||

// IndexAt returns the index (i.e. the sample #) at a specified TimeStamp.
func (t *Chunk) IndexAt(ts TimeStamp) int64 {
	return int64(t.RangeFromStart(ts).Span() / t.Period())
}

// ByteIndexAt returns the byte index representing the location of the sample data at the provided TimeStamp.
func (t *Chunk) ByteIndexAt(ts TimeStamp) int64 {
	return t.IndexAt(ts) * t.SampleSize()
}

// |||| VALUE ACCESS ||||

// ValueAt returns the value at the provided TimeStamp. The returned value will have the same type as Chunk.DataType.
func (t *Chunk) ValueAt(ts TimeStamp) interface{} {
	byteI := t.ByteIndexAt(ts)
	return convertBytes(t.ReadSlice(byteI, byteI+t.SampleSize()), t.DataType)
}

// ValuesInRange returns all values in the provided TimeRange.
// The returned value will be a slice with the same element type as Chunk.DataType
// If the provided TimeRange exceeds the bounds of the chunk, will not return an error,
// and will instead return an empty slice.
func (t *Chunk) ValuesInRange(rng TimeRange) interface{} {
	capped := t.capRange(rng)
	startByteI, endByteI := t.ByteIndexAt(capped.Start()), t.ByteIndexAt(capped.End())
	return convertBytes(t.ReadSlice(startByteI, endByteI), t.DataType)
}

// AllValues returns a slice representing all the values in the chunk. The elements in the slice will have
// the same type as Chunk.DataType.
// WARNING: This is an expensive operation as it requires converting bytes to concrete data types. Be wary
// when calling this with large chunks.
func (t *Chunk) AllValues() interface{} {
	return convertBytes(t.Bytes(), t.DataType)
}

// |||| MODIFICATION ||||

// RemoveFromStart removes all values from the Start of the chunk to the specified timestamp. Modifies the Start
// of the chunk to account for the removed values.
//
// Will panic if the provided TimeStamp exceeds the Range of the chunk.
func (t *Chunk) RemoveFromStart(ts TimeStamp) {
	t.Splice(t.ByteIndexAt(t.Start()), t.ByteIndexAt(ts))
	t.start = ts
}

// RemoveFromEnd removes all values from the End of the chunk to the specified timestamp. Modifies the End of the
// chunk to account for the removed values.
//
// Will panic if the provided TimeStamp exceeds the Range of the chunk.
func (t *Chunk) RemoveFromEnd(ts TimeStamp) {
	t.Splice(t.ByteIndexAt(ts), t.ByteIndexAt(t.End()))

}

// |||| OVERLAP ||||

// Overlap returns a new ChunkOverlap representing the overlap between Chunk and a provided Chunk cChunk
// The provided Chunk will be considered the dest Chunk.
func (t *Chunk) Overlap(cChunk *Chunk) ChunkOverlap {
	return ChunkOverlap{source: t, dest: cChunk}
}

// |||| UTILITIES ||||

// lastValueTS returns the TimeStamp of the last value in the chunk.
// Because chunk.End is exclusive, retrieving the value at that TimeStamp would
// result in an error, so we need to find that last inclusive TimeStamp.
func (t *Chunk) lastValueTS() TimeStamp {
	return t.End().Add(-1 * t.Period())
}

func (t *Chunk) capRange(rng TimeRange) TimeRange {
	return NewTimeRange(t.capStamp(rng.Start()), t.capStamp(rng.End()))
}

func (t *Chunk) capStamp(ts TimeStamp) TimeStamp {
	if ts > t.lastValueTS() {
		return t.End()
	}
	if ts < t.Start() {
		return t.Start()
	}
	return ts
}

func sampleSize(dt DataType) int64 {
	switch dt {
	case DataTypeFloat64:
		return 8
	case DataTypeFloat32:
		return 4
	}
	panic("t chunk has unknown data type")
}
