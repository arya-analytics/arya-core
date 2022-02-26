package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local", func() {
	var (
		localSvc      chanchunk.Local
		items         []interface{}
		channelConfig *models.ChannelConfig
		node          *models.Node
		rng           *models.Range
		ccr           *models.ChannelChunkReplica
		rr            *models.RangeReplica
		cc            *models.ChannelChunk
	)
	BeforeEach(func() {
		localSvc = chanchunk.NewServiceLocalStorage(store)
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rng = &models.Range{ID: uuid.New()}
		rr = &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
		cc = &models.ChannelChunk{ID: uuid.New(), RangeID: rng.ID, ChannelConfigID: channelConfig.ID}
		ccr = &models.ChannelChunkReplica{RangeReplicaID: rr.ID, ChannelChunkID: cc.ID, Telem: telem.NewBulk([]byte("randomdata"))}
		items = []interface{}{node, channelConfig, rng, rr, cc}
	})
	JustBeforeEach(func() {
		for _, m := range items {
			Expect(store.NewCreate().Model(m).Exec(ctx)).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, m := range items {
			Expect(store.NewDelete().Model(m).WherePK(model.NewReflect(m).PK().Raw()).Exec(ctx)).To(BeNil())
		}
	})
	Describe("Channel Chunk Replica", func() {
		JustBeforeEach(func() {
			Expect(localSvc.Create(ctx, ccr)).To(BeNil())
		})
		JustAfterEach(func() {
			ccrOpts := chanchunk.LocalDeleteOpts{
				PKC: model.NewPKChain([]uuid.UUID{cc.ID}),
			}
			Expect(localSvc.Delete(ctx, ccrOpts)).To(BeNil())
		})
		It("Should create the replica correctly", func() {
			resCCR := &models.ChannelChunkReplica{}
			opts := chanchunk.LocalRetrieveOpts{
				PKC: model.NewPKChain([]uuid.UUID{ccr.ID}),
			}
			err := localSvc.Retrieve(ctx, resCCR, opts)
			Expect(err).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
			Expect(resCCR.Telem.Bytes()).To(Equal([]byte("randomdata")))
		})
		It("Should omit bulk data when the field is not specified", func() {
			resCCR := &models.ChannelChunkReplica{}
			opts := chanchunk.LocalRetrieveOpts{
				PKC:    model.NewPKChain([]uuid.UUID{ccr.ID}),
				Fields: []string{"ID", "RangeReplicaID", "ChannelChunkID"},
			}
			err := localSvc.Retrieve(ctx, resCCR, opts)
			Expect(err).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
			Expect(resCCR.Telem).To(BeNil())
		})
		It("Should bulk update the replica correctly", func() {
			newRR := &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
			Expect(store.NewCreate().Model(newRR).Exec(ctx)).To(BeNil())
			updateCCR := []*models.ChannelChunkReplica{{ID: ccr.ID, RangeReplicaID: newRR.ID}}
			opts := chanchunk.LocalUpdateOpts{
				Bulk:   true,
				Fields: []string{"RangeReplicaID"},
			}
			Expect(localSvc.Update(ctx, &updateCCR, opts)).To(BeNil())
			resCCR := &models.ChannelChunkReplica{}
			rOpts := chanchunk.LocalRetrieveOpts{WhereFields: model.WhereFields{"RangeReplicaID": newRR.ID}}
			Expect(localSvc.Retrieve(ctx, resCCR, rOpts)).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
		})
	})
	Describe("Retrieve Range Replicas", func() {
		It("Should retrieve the node information correctly", func() {
			resRR := &models.RangeReplica{}
			err := localSvc.RetrieveRangeReplica(ctx, resRR, model.NewPKChain([]uuid.UUID{rr.ID}))
			Expect(err).To(BeNil())
			Expect(resRR.Node.ID).To(Equal(1))
		})
	})
})
