package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	BeforeEach(migrate)
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct error type", func() {
				somePKThatDoesntExist := 136987
				m := &storage.ChannelConfig{}
				err := dummyEngine.NewRetrieve(dummyAdapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(dummyCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeItemNotFound))
			})
		})
	})
})
