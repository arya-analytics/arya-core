package telem

import (
	"reflect"
)

// ChunkOverlap queries and manages the overlap between two different Chunk of telemetry.
// Provides functionality for assessing the overlap, and resolving the overlap by modifying the underlying Chunk.
type ChunkOverlap struct {
	source *Chunk
	dest   *Chunk
}

// OverlapType classifies an overlap into several distinct types.
type OverlapType int

// The following examples show different types of overlaps. Each example contains a source slice and a destination slice
// (in that order). The values in the slices represent timestamps of samples.
//go:generate stringer -type=OverlapType
const (
	// OverlapTypeNoneOrInvalid is an overlap that is either nonexistent (gap or contiguous) or an invalid overlap,
	// such as conflicting data rates or types.
	OverlapTypeNoneOrInvalid OverlapType = iota
	// OverlapTypeRightPartial indicates an overlap where the source Chunk ends after the
	// destination chunk. Ex:
	//
	// [3,4,5], [1,2,3]
	//
	OverlapTypeRightPartial
	// OverlapTypeLeftPartial indicates an overlap where the source Chunk starts before the dest chunk. Ex:
	//
	// [1,2,3], [3,4,5]
	//
	OverlapTypeLeftPartial
	// OverlapTypeSourceConsume indicates an overlap where the source Chunk 'consumes' or envelops the dest chunk. Ex:
	//
	// [1,2,3,4,5,6], [3,4,5]
	OverlapTypeSourceConsume
	// OverlapTypeDestConsume indicates an overlap where the dest Chunk 'consumes' or envelops the source chunk. Ex:
	//
	// [3,4,5], [1,2,3,4,5,6]
	//
	OverlapTypeDestConsume
	// OverlapTypeDuplicate indicates an overlap where the source and dest chunks occupy the same TimeRange. Ex:
	//
	// [1,2,3], [1,2,3]
	OverlapTypeDuplicate
)

/// |||| GENERIC INFO ||||

// TimeStampExp returns a TimeRange representing the range of the ChunkOverlap. This might be an invalid or negative range
// if the ChunkOverlap is invalid.
func (o ChunkOverlap) Range() TimeRange {
	r, _ := o.baseRange()
	return r
}

// Type returns the OverlapType of the ChunkOverlap. See the OverlapType documentation for more on this.
func (o ChunkOverlap) Type() OverlapType {
	if !o.IsValid() {
		return OverlapTypeNoneOrInvalid
	}
	if o.source.Start() < o.dest.Start() && o.source.End() > o.dest.End() {
		return OverlapTypeSourceConsume
	}
	if o.source.Start() > o.dest.Start() && o.source.End() < o.dest.End() {
		return OverlapTypeDestConsume
	}
	if o.source.Start() == o.dest.Start() && o.source.Span() == o.dest.Span() {
		return OverlapTypeDuplicate
	}
	if o.source.Start() < o.dest.Start() {
		return OverlapTypeLeftPartial
	}
	if o.source.End() > o.dest.End() {
		return OverlapTypeRightPartial
	}
	panic("could not determine overlap type")
}

// |||| VALIDATION ||||

// IsValid checks if the ChunkOverlap is valid. This falls on two criteria:
//
// 1. The chunks actually overlap i.e. they occupy a common TimeRange.
// 2. The chunks are compatible (see ChunksCompatible)
//
func (o ChunkOverlap) IsValid() bool {
	_, valid := o.baseRange()
	return valid && o.ChunksCompatible()
}

// IsUniform checks if the ChunkOverlap is 'uniform' i.e. the telemetry residing in source is the same as the telemetry
// residing in dest for the overlap TimeStampExp.
//
// This essentially checks if the data between two chunks is contiguous or not.
//
// WARNING: This calls reflect.DeepEqual internally, which is expensive and can cause perf issues.
func (o ChunkOverlap) IsUniform() bool {
	if !o.IsValid() {
		return false
	}
	return reflect.DeepEqual(o.SourceValues(), o.DestValues())
}

// ChunksCompatible checks if the source and dest Chunk are 'compatible' with one another. They are compatible
// if they meet the following criteria:
//
// 1. The Chunk.DataRate are the same.
// 2. the Chunk.DataType are the same.
//
func (o ChunkOverlap) ChunksCompatible() bool {
	return o.dest.DataRate == o.source.DataRate && o.dest.DataType == o.source.DataType
}

// |||| VALUE ACCESS ||||

// SourceValues returns the telemetry in the source Chunk within the ChunkOverlap TimeStampExp.
func (o ChunkOverlap) SourceValues() interface{} {
	return o.source.ValuesInRange(o.Range())
}

// DestValues returns the telemetry in the dest Chunk within the ChunkOverlap TimeStampExp.
func (o ChunkOverlap) DestValues() interface{} {
	return o.dest.ValuesInRange(o.Range())
}

// |||| MODIFICATION ||||

// RemoveFromSource removes the overlapping values from the source chunk.
//
// If the overlap is of OverlapType OverlapTypeDuplicate, will erase the entire source chunk.
//
// WARNING: This operation will panic if the ChunkOverlap is of OverlapType OverlapTypeSourceConsume or
// OverlapTypeDestConsume. For removing consumed ChunkOverlap, see RemoveFromConsumed.
func (o ChunkOverlap) RemoveFromSource() error {
	return o.removeFrom(o.source)
}

// RemoveFromDest is the same as RemoveFromSource, but removes from the dest chunk.
func (o ChunkOverlap) RemoveFromDest() error {
	return o.removeFrom(o.dest)
}

// RemoveFromConsumed removes the overlapping values from the consumed chunk i.e. erases the consumed chunks data.
//
// WARNING: This operation will panic if the ChunkOverlap is not of OverlapType OverlapTypeSourceConsume or
// OverlapTypeDestConsume.
func (o ChunkOverlap) RemoveFromConsumed() error {
	var c *Chunk
	if o.Type() == OverlapTypeDestConsume {
		c = o.source
	} else if o.Type() == OverlapTypeSourceConsume {
		c = o.dest
	}
	if c == nil {
		panic("can't call remove from consumed on non consume overlap")
	}
	_, err := c.Write([]byte{})
	return err
}

func (o ChunkOverlap) removeFrom(c *Chunk) error {
	if o.Type() == OverlapTypeDestConsume || o.Type() == OverlapTypeSourceConsume {
		panic("can't use removeFrom on consuming overlaps")
	} else if o.Type() == OverlapTypeNoneOrInvalid {
		panic("no overlap to remove!")
	}
	if o.Type() == OverlapTypeDuplicate {
		_, err := c.Write([]byte{})
		return err
	}
	if o.Range().Start() == c.Start() {
		c.RemoveFromStart(o.Range().End())
		return nil
	}
	c.RemoveFromEnd(o.Range().Start())
	return nil
}

func (o ChunkOverlap) baseRange() (TimeRange, bool) {
	return o.source.Range().Overlap(o.dest.Range())
}
