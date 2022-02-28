package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	var channelChunk *models.ChannelChunkReplica
	BeforeEach(func() {
		channelChunk = &models.ChannelChunkReplica{
			ID:    uuid.New(),
			Telem: telem.NewChunkData([]byte("randomstring")),
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
		mockModelTwo := &models.ChannelChunkReplica{}
		err := engine.NewRetrieve(adapter).Model(mockModelTwo).WherePK(channelChunk.ID).
			Exec(ctx)
		Expect(err).To(BeNil())
		Expect(mockModelTwo.Telem).ToNot(BeNil())
		Expect(mockModelTwo.Telem.Bytes()).To(Equal([]byte("randomstring")))
	})
})
