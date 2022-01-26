package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	var mockCreate *storage.ChannelChunk
	var err error
	BeforeEach(func() {
		mockCreate = &storage.ChannelChunk{
			ID:   uuid.New(),
			Data: mock.NewObject([]byte("randomstring")),
		}
		err = mockEngine.NewCreate(mockAdapter).Model(mockCreate).Exec(mockCtx)
	})
	AfterEach(func() {
		err = mockEngine.NewDelete(mockAdapter).Model(mockCreate).WherePK(mockCreate.ID).Exec(
			mockCtx)
	})
	It("Should create without errutil", func() {
		Expect(err).To(BeNil())
	})
	It("Should be able to be re-queried after creation", func() {
		mockModelTwo := &storage.ChannelChunk{}
		err := mockEngine.NewRetrieve(mockAdapter).Model(mockModelTwo).WherePK(
			mockCreate.ID).
			Exec(mockCtx)
		Expect(err).To(BeNil())
		Expect(mockModelTwo.Data).ToNot(BeNil())
		b := make([]byte, mockCreate.Data.Size())
		_, err = mockModelTwo.Data.Read(b)
		Expect(err.Error()).To(Equal("EOF"))
		Expect(b).To(Equal([]byte("randomstring")))
	})
})
