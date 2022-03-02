package telem

import (
	"encoding/binary"
	"io"
)

// ChunkData stores telemetry in its binary form, and exposes io.Reader and io.Writer interfaces.
type ChunkData struct {
	pos  int
	data []byte
}

// NewChunkData instantiates and returns a new ChunkData. The data argument
// should generally be left blank and should be a byte slice with the amount of data needing to be read.
func NewChunkData(data []byte) *ChunkData {
	return &ChunkData{data: data}
}

// || READ ||

// Read implements io.Reader.Read
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
	return i, nil
}

// Size returns the size of the ChunkData in bytes.
func (cd *ChunkData) Size() int64 {
	return int64(len(cd.data))
}

// Bytes returns the ChunkData as a byte array.
func (cd *ChunkData) Bytes() []byte {
	return cd.data
}

// Done returns true of the entire ChunkData has been read.
func (cd *ChunkData) Done() bool {
	return cd.pos == len(cd.data)
}

// Reset resets the underlying cursor to the beginning of the ChunkData.
func (cd *ChunkData) Reset() {
	cd.pos = 0
}

// ReadSlice reads a slice of bytes from one index to another.
func (cd *ChunkData) ReadSlice(from int64, to int64) []byte {
	return cd.data[from:to]
}

// || WRITE ||

// ReadFrom implements io.Writer.
func (cd *ChunkData) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(cd.data)
	if err != io.EOF {
		return int64(n), err
	}
	return int64(n), nil
}

// Write implements io.Writer.
func (cd *ChunkData) Write(p []byte) (n int, err error) {
	cd.data = make([]byte, len(p))
	for i, b := range p {
		cd.data[i] = b
	}
	cd.Reset()
	return len(p), nil
}

// WriteData writes an arbitrary interface to the ChunkData.
func (cd *ChunkData) WriteData(data interface{}) error {
	return binary.Write(cd, ByteOrder(), data)
}

// || MODIFY ||

// Splice splices and removes the data in the chunk from the start to the end indexes provided.
// from index is inclusive, to index is exclusive.
func (cd *ChunkData) Splice(from int64, to int64) {
	cd.data = append(cd.data[:from], cd.data[to:]...)
}
