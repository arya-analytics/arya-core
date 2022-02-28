package telem

import "reflect"

type ChunkOverlap struct {
	source *Chunk
	dest   *Chunk
}

func (o ChunkOverlap) Range() TimeRange {
	r, _ := o.baseRange()
	return r
}

func (o ChunkOverlap) IsValid() bool {
	_, valid := o.baseRange()
	return valid && o.chunksCompatible()
}

func (o ChunkOverlap) IsUniform() bool {
	if !o.IsValid() {
		return false
	}
	return reflect.DeepEqual(o.SourceValues(), o.DestValues())
}

func (o ChunkOverlap) SourceValues() interface{} {
	return o.source.ValuesInRange(o.Range())
}

func (o ChunkOverlap) DestValues() interface{} {
	return o.dest.ValuesInRange(o.Range())
}

func (o ChunkOverlap) chunksCompatible() bool {
	return o.dest.dataRate == o.source.dataRate && o.dest.dataType == o.source.dataType
}

func (o ChunkOverlap) baseRange() (TimeRange, bool) {
	return o.source.Range().Overlap(o.dest.Range())
}
