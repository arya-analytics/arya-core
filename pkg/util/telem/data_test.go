package telem_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Telem", func() {
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
	It("Should read bytes out of the chunk data", func() {
		var p = make([]byte, 3)
		cd.Write([]byte{1, 2, 3})
		n, err := cd.Read(p)
		Expect(n).To(Equal(3))
		Expect(err).To(BeNil())
		Expect(p[1]).To(Equal(uint8(2)))
		Expect(cd.Done()).To(BeTrue())
		cd.Reset()
	})
	It("Should be able to read bytes multiple times", func() {
		var pOne = make([]byte, 3)
		cd.Write([]byte{1, 2, 3, 4, 5})
		n, err := cd.Read(pOne)
		Expect(n).To(Equal(3))
		Expect(err).To(BeNil())
		Expect(pOne[1]).To(Equal(uint8(2)))

		Expect(cd.Done()).To(BeFalse())

		var pTwo = make([]byte, 2)

		n, err = cd.Read(pTwo)

		Expect(n).To(Equal(2))
		Expect(pTwo).To(Equal([]byte{4, 5}))
		Expect(err).To(BeNil())

		Expect(cd.Done()).To(BeTrue())
	})
	It("Should read bytes into another chunk data", func() {
		cd.Write([]byte{1, 2, 3, 4, 5})
		cdTwo := telem.NewChunkData(make([]byte, 5))
		n, err := cdTwo.ReadFrom(cd)
		Expect(err).To(BeNil())
		Expect(n).To(Equal(int64(5)))
	})
})
