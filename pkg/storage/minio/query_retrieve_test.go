package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryRetrieve", func() {
	var (
		channelChunk *models.ChannelChunkReplica
		bytes        []byte
	)
	BeforeEach(func() {
		bytes = []byte("randomstring")
		channelChunk = &models.ChannelChunkReplica{
			ID:    uuid.New(),
			Telem: telem.NewBulk(bytes),
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
	Describe("Standard Usage", func() {
		Describe("Retrieve an item", func() {
			It("Should retrieve the correct item", func() {
				resChannelChunk := &models.ChannelChunkReplica{}
				err := engine.NewRetrieve(adapter).Model(resChannelChunk).WherePK(channelChunk.ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resChannelChunk.Telem).ToNot(BeNil())
				Expect(resChannelChunk.Telem.Bytes()).To(Equal([]byte("randomstring")))
			})
		})
		Describe("Retrieve multiple items", func() {
			var channelChunkTwo *models.ChannelChunkReplica
			BeforeEach(func() {
				channelChunkTwo = &models.ChannelChunkReplica{
					ID:    uuid.New(),
					Telem: telem.NewBulk([]byte("model two")),
				}
			})
			JustBeforeEach(func() {
				err := engine.NewCreate(adapter).Model(channelChunkTwo).Exec(
					ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct items", func() {
				var models []*models.ChannelChunkReplica
				err := engine.NewRetrieve(adapter).Model(&models).WherePKs([]uuid.
					UUID{channelChunk.ID, channelChunkTwo.ID}).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(models).To(HaveLen(2))
				Expect([]uuid.UUID{channelChunk.ID, channelChunkTwo.ID}).To(
					ContainElements(models[0].ID, models[1].ID))
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesnt exist", func() {
			It("Should return the correct error type", func() {
				somePKThatDoesntExist := uuid.New()
				err := engine.NewRetrieve(adapter).Model(channelChunk).WherePK(
					somePKThatDoesntExist).Exec(ctx)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
