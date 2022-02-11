package telem

import "bytes"

type Bulk struct {
	*bytes.Buffer
}

func NewBulk(buf []byte) *Bulk {
	return &Bulk{
		Buffer: bytes.NewBuffer(buf),
	}
}

func (b *Bulk) Size() int64 {
	return int64(b.Len())
}
