package telem

import (
	"reflect"
)

type ChunkOverlap struct {
	source *Chunk
	dest   *Chunk
}

type OverlapType int

const (
	OverlapTypeNoneOrInvalid OverlapType = iota
	OverlapTypePartial
	OverlapTypeSourceConsume
	OverlapTypeDestConsume
	OverlapTypeDuplicate
)

/// |||| GENERIC INFO ||||

func (o ChunkOverlap) Range() TimeRange {
	r, _ := o.baseRange()
	return r
}

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
	return OverlapTypePartial
}

// |||| VALIDATION ||||

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

// |||| VALUE ACCESS ||||

func (o ChunkOverlap) SourceValues() interface{} {
	return o.source.ValuesInRange(o.Range())
}

func (o ChunkOverlap) DestValues() interface{} {
	return o.dest.ValuesInRange(o.Range())
}

// |||| MODIFICATION ||||

func (o ChunkOverlap) RemoveFromSource() error {
	return o.removeFrom(o.source)
}

func (o ChunkOverlap) RemoveFromDest() error {
	return o.removeFrom(o.dest)
}

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

func (o ChunkOverlap) chunksCompatible() bool {
	return o.dest.dataRate == o.source.dataRate && o.dest.dataType == o.source.dataType
}

func (o ChunkOverlap) baseRange() (TimeRange, bool) {
	return o.source.Range().Overlap(o.dest.Range())
}
