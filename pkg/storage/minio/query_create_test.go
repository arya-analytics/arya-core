package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	var mockModel *storage.ChannelChunk
	It("Should create without error", func() {
		mockModel = &storage.ChannelChunk{
			ID:   uuid.New(),
			Data: mock.NewObject([]byte("randomstring")),
		}
		err := mockEngine.NewCreate(mockAdapter).Model(mockModel).Exec(mockCtx)
		Expect(err).To(BeNil())
	})
	It("Should be able to be re-queried after creation", func() {
		mockModelTwo := &storage.ChannelChunk{}
		err := mockEngine.NewRetrieve(mockAdapter).Model(mockModelTwo).WherePK(
			mockModel.ID).
			Exec(mockCtx)
		Expect(err).To(BeNil())
		Expect(mockModelTwo.Data).ToNot(BeNil())
		b := make([]byte, mockModel.Data.Size())
		_, err = mockModelTwo.Data.Read(b)
		Expect(err.Error()).To(Equal("EOF"))
		Expect(b).To(Equal([]byte("randomstring")))
	})
})
