package telem

// |||| DATA TYPE ||||

type DataType int

//go:generate stringer -type=DataType
const (
	DataTypeFloat64 DataType = iota + 1
	DataTypeFloat32
)

// |||| CHUNK ||||

type Chunk struct {
	start    TimeStamp
	DataType DataType
	DataRate DataRate
	*ChunkData
}

func NewChunk(start TimeStamp, dataType DataType, dataRate DataRate, data *ChunkData) *Chunk {
	return &Chunk{start: start, DataType: dataType, DataRate: dataRate, ChunkData: data}
}

// |||| STATIC ATTRIBUTES ||||

func (t *Chunk) SampleSize() int64 {
	return sampleSize(t.DataType)
}

func (t *Chunk) Period() TimeSpan {
	return t.DataRate.Period()
}

// |||| SIZING ||||

func (t *Chunk) Len() int64 {
	return t.Size() / t.SampleSize()
}

// |||| TIMING ||||

func (t *Chunk) RangeSinceStart(ts TimeStamp) TimeRange {
	return NewTimeRange(t.start, ts)
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
	return TimeSpan(t.Len() * int64(t.DataRate.Period()))
}

// |||| INDEXING ||||

func (t *Chunk) IndexAtTS(ts TimeStamp) int64 {
	return int64(t.RangeSinceStart(ts).Span() / t.Period())
}

func (t *Chunk) ByteIndexAtTS(ts TimeStamp) int64 {
	return t.IndexAtTS(ts) * t.SampleSize()
}

// |||| VALUE ACCESS ||||

func (t *Chunk) ValueAtTS(ts TimeStamp) interface{} {
	byteI := t.ByteIndexAtTS(ts)
	return convertBytes(t.ReadSlice(byteI, byteI+t.SampleSize()), t.DataType)
}

func (t *Chunk) ValuesInRange(rng TimeRange) interface{} {
	capped := t.capRange(rng)
	startByteI, endByteI := t.ByteIndexAtTS(capped.Start()), t.ByteIndexAtTS(capped.End())
	return convertBytes(t.ReadSlice(startByteI, endByteI), t.DataType)
}

func (t *Chunk) AllValues() interface{} {
	return convertBytes(t.Bytes(), t.DataType)
}

// |||| MODIFICATION ||||

func (t *Chunk) RemoveFromStart(ts TimeStamp) {
	t.Splice(t.ByteIndexAtTS(t.Start()), t.ByteIndexAtTS(ts))
	t.start = ts
}

func (t *Chunk) RemoveFromEnd(ts TimeStamp) {
	t.Splice(t.ByteIndexAtTS(ts), t.ByteIndexAtTS(t.End()))

}

// |||| OVERLAP ||||

func (t *Chunk) Overlap(cChunk *Chunk) ChunkOverlap {
	return ChunkOverlap{source: t, dest: cChunk}
}

// |||| UTILITIES ||||

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
