package telem_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {
	var (
		cd *telem.ChunkData
	)
	BeforeEach(func() {
		cd = telem.NewChunkData([]byte{})
	})
	It("Should create a new empty ChunkData", func() {
		Expect(cd.Size()).To(Equal(int64(0)))
	})
	It("Should write bytes to the chunk data", func() {
		n, err := cd.Write([]byte{1, 2, 3})
		Expect(n).To(Equal(3))
		Expect(err).To(BeNil())
	})
	It("Should read bytes out fo the chunk data", func() {
		var p = make([]byte, 3)
		cd.Write([]byte{1, 2, 3})
		n, err := cd.Read(p)
		Expect(n).To(Equal(3))
		Expect(err).To(BeNil())
		Expect(p[1]).To(Equal(uint8(2)))
		Expect(cd.Done()).To(Equal(true))
		cd.Reset()
	})

})
