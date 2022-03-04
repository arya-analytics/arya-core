package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
			err := engine.NewCreate().Model(channelChunk).Exec(ctx)
			Expect(err).To(BeNil())
		})
		Describe("Delete an item", func() {
			JustBeforeEach(func() {
				err := engine.NewDelete().Model(channelChunk).WherePK(
					channelChunk.ID).Exec(
					ctx)
				Expect(err).To(BeNil())
			})
			It("Should not be able to be re-queried after del", func() {
				rErr := engine.NewRetrieve().Model(channelChunk).WherePK(
					channelChunk.ID).Exec(ctx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
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
				cErr := engine.NewCreate().Model(channelChunkTwo).Exec(ctx)
				Expect(cErr).To(BeNil())
				pks := []uuid.UUID{channelChunk.ID, channelChunkTwo.ID}
				dErr := engine.NewDelete().Model(channelChunkTwo).WherePKs(pks).Exec(ctx)
				Expect(dErr).To(BeNil())
			})
			It("Should not be able to be re-queried after del", func() {
				var models []*models.ChannelChunkReplica
				e := engine.NewRetrieve().Model(&models).WherePKs(
					[]uuid.UUID{channelChunkTwo.ID, channelChunk.ID}).Exec(ctx)
				Expect(e).ToNot(BeNil())
				Expect(e.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
			})
		})
	})
})
