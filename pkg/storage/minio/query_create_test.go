package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	var channelChunk *storage.ChannelChunkReplica
	BeforeEach(func() {
		channelChunk = &storage.ChannelChunkReplica{
			ID:    uuid.New(),
			Telem: mock.NewObject([]byte("randomstring")),
		}
	})
	JustBeforeEach(func() {
		err := engine.NewCreate(adapter).Model(channelChunk).Exec(ctx)
		Expect(err).To(BeNil())
	})
	AfterEach(func() {
		err := engine.NewDelete(adapter).Model(channelChunk).WherePK(channelChunk.
			ID).Exec(ctx)
		Expect(err).To(BeNil())
	})
	It("Should be created correctly", func() {
		mockModelTwo := &storage.ChannelChunkReplica{}
		err := engine.NewRetrieve(adapter).Model(mockModelTwo).WherePK(channelChunk.ID).
			Exec(ctx)
		Expect(err).To(BeNil())
		Expect(mockModelTwo.Telem).ToNot(BeNil())
		b := make([]byte, channelChunk.Telem.Size())
		_, err = mockModelTwo.Telem.Read(b)
		Expect(err.Error()).To(Equal("EOF"))
		Expect(b).To(Equal([]byte("randomstring")))
	})
})
