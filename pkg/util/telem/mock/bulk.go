package mock

import (
	"encoding/binary"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"math/rand"
)

func TelemBulkPopulateRandomFloat64(tb *telem.Bulk, qty int) {
	var w64 []float64
	for i := 0; i < qty; i++ {
		w64 = append(w64, rand.Float64())
	}
	if err := writeSlice(tb, w64); err != nil {
		panic(err)
	}
}

func TelemBulkPopulateRandomFloat32(tb *telem.Bulk, qty int) {
	var w32 []float32
	for i := 0; i < qty; i++ {
		w32 = append(w32, rand.Float32())
	}
	if err := writeSlice(tb, w32); err != nil {
		panic(err)
	}
}

func writeSlice(tb *telem.Bulk, slice interface{}) error {
	return binary.Write(tb, binary.LittleEndian, slice)
}
