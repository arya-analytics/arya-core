package rpc_test

import (
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type BulkTelemModel struct {
	Telem *telem.ChunkData `model:"role:telemChunkData"`
}

type BytesTelemModel struct {
	Telem []byte `model:"role:telemChunkData"`
}

var _ = Describe("ModelExchange", func() {
	Describe("FieldHandlerTelemBulk", func() {
		Context("Standard usage", func() {
			blk := telem.NewChunkData([]byte{})
			mock.PopulateRandFloat64(blk, 100)
			It("Should exchange to dest correctly", func() {
				bulkModel := &BulkTelemModel{
					Telem: blk,
				}
				bytesModel := &BytesTelemModel{}
				me := rpc.NewModelExchange(bulkModel, bytesModel)
				me.ToDest()
				Expect(int64(len(bytesModel.Telem))).To(Equal(blk.Size()))
			})
			It("Should exchange to source correctly", func() {
				bulkModel := &BulkTelemModel{}
				bytesModel := &BytesTelemModel{
					Telem: blk.Bytes(),
				}
				me := rpc.NewModelExchange(bulkModel, bytesModel)
				me.ToSource()
				Expect(bulkModel.Telem.Size()).To(Equal(blk.Size()))
			})
		})
	})
})
