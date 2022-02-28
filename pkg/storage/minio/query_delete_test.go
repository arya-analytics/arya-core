package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryDelete", func() {
	var channelChunk *models.ChannelChunkReplica
	Describe("Standard Usage", func() {
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
		Describe("Delete an item", func() {
			JustBeforeEach(func() {
				err := engine.NewDelete(adapter).Model(channelChunk).WherePK(
					channelChunk.ID).Exec(
					ctx)
				Expect(err).To(BeNil())
			})
			It("Should not be able to be re-queried after delete", func() {
				rErr := engine.NewRetrieve(adapter).Model(channelChunk).WherePK(
					channelChunk.ID).Exec(ctx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
		Describe("Delete multiple Items", func() {
			var channelChunkTwo *models.ChannelChunkReplica
			BeforeEach(func() {
				channelChunkTwo = &models.ChannelChunkReplica{
					ID:    uuid.New(),
					Telem: telem.NewChunkData([]byte("mock bytes")),
				}
			})
			JustBeforeEach(func() {
				cErr := engine.NewCreate(adapter).Model(channelChunkTwo).Exec(ctx)
				Expect(cErr).To(BeNil())
				pks := []uuid.UUID{channelChunk.ID, channelChunkTwo.ID}
				dErr := engine.NewDelete(adapter).Model(channelChunkTwo).WherePKs(pks).Exec(ctx)
				Expect(dErr).To(BeNil())
			})
			It("Should not be able to be re-queried after delete", func() {
				var models []*models.ChannelChunkReplica
				e := engine.NewRetrieve(adapter).Model(&models).WherePKs(
					[]uuid.UUID{channelChunkTwo.ID, channelChunk.ID}).Exec(ctx)
				Expect(e).ToNot(BeNil())
				Expect(e.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
		})
	})
})
