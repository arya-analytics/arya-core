package telem

import (
	"encoding/binary"
	"io"
)

type ChunkData struct {
	pos  int
	data []byte
}

func NewChunkData(data []byte) *ChunkData {
	return &ChunkData{data: data}
}

// || READ ||

func (cd *ChunkData) Read(p []byte) (n int, err error) {
	if cd.Done() {
		return 0, io.EOF
	}
	i := 0
	for ; i < len(cd.data)-cd.pos; i++ {
		if i > len(p)-1 {
			break
		}
		p[i] = cd.data[i+cd.pos]
	}
	cd.pos += i
	return len(p), nil
}

func (cd *ChunkData) Size() int64 {
	return int64(len(cd.data))
}

func (cd *ChunkData) Bytes() []byte {
	return cd.data
}

func (cd *ChunkData) Done() bool {
	return cd.pos == len(cd.data)
}

func (cd *ChunkData) Reset() {
	cd.pos = 0
}

func (cd *ChunkData) ReadSlice(from int64, to int64) []byte {
	return cd.data[from:to]
}

// || WRITE ||

func (cd *ChunkData) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(cd.data)
	if err != io.EOF {
		return int64(n), err
	}
	return int64(n), nil
}

func (cd *ChunkData) Write(p []byte) (n int, err error) {
	cd.data = make([]byte, len(p))
	for i, b := range p {
		cd.data[i] = b
	}
	cd.Reset()
	return len(p), nil
}

func (cd *ChunkData) WriteData(data interface{}) error {
	return binary.Write(cd, binary.BigEndian, data)
}

// || MODIFY ||

func (cd *ChunkData) Splice(from int64, to int64) {
	cd.data = append(cd.data[:from], cd.data[to:]...)
}
