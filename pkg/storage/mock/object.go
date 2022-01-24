package mock

import (
	"bytes"
	"io"
)

type Object struct {
	b []byte
	r io.Reader
}

func NewObject(b []byte) *Object {
	return &Object{
		b: b,
		r: bytes.NewReader(b),
	}
}

func (mo *Object) Read(b []byte) (n int, err error) {
	return mo.r.Read(b)
}
func (mo *Object) Size() int64 {
	return int64(len(mo.b))
}
