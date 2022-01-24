package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	It("Should create without error", func() {
		mockModel := &storage.ChannelChunk{
			ID:   uuid.New(),
			Data: mock.NewObject([]byte("randomstring")),
		}
		err := mockEngine.NewCreate(mockAdapter).Model(mockModel).Exec(mockCtx)
		Expect(err).To(BeNil())
	})
})
