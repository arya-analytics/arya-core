package mock

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"math/rand"
)

func PopulateRandFloat64(cd *telem.ChunkData, qty int) {
	var w64 []float64
	for i := 0; i < qty; i++ {
		w64 = append(w64, rand.Float64())
	}
	if err := cd.WriteData(w64); err != nil {
		panic(err)
	}
}

func PopulateRandFloat32(cd *telem.ChunkData, qty int) {
	var w32 []float32
	for i := 0; i < qty; i++ {
		w32 = append(w32, rand.Float32())
	}
	if err := cd.WriteData(w32); err != nil {
		panic(err)
	}
}

func PopulateRand(cd *telem.ChunkData, dataType telem.DataType, qty int) {
	switch dataType {
	case telem.DataTypeFloat32:
		PopulateRandFloat32(cd, qty)
	case telem.DataTypeFloat64:
		PopulateRandFloat64(cd, qty)
	default:
		panic("can't populate unknown data type with random values")
	}
}

func spanToSampleQTY(span telem.TimeSpan, dataRate telem.DataRate) int {
	return int(span / dataRate.Period())
}

func ChunkSet(
	chunkQty int,
	start telem.TimeStamp,
	dataType telem.DataType,
	dataRate telem.DataRate,
	chunkSpan telem.TimeSpan,
	gap telem.TimeSpan,
) (chunks []*telem.Chunk) {
	sampleQty := spanToSampleQTY(chunkSpan, dataRate)
	for i := 0; i < chunkQty; i++ {
		cd := telem.NewChunkData([]byte{})
		PopulateRand(cd, dataType, sampleQty)
		var startTS telem.TimeStamp
		if i == 0 {
			startTS = start
		} else {
			startTS = chunks[i-1].End().Add(gap)
		}
		chunks = append(chunks, telem.NewChunk(startTS, dataType, dataRate, cd))
	}
	return chunks
}
