package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Persist", func() {
	Describe("PersistCluster", func() {
		var (
			newRng  *models.Range
			newRR   *models.RangeReplica
			chanCfg *models.ChannelConfig
			p       rng.Persist
			items   []interface{}
		)
		BeforeEach(func() {
			if clust == nil {
				var err error
				clust, err = mock.New(ctx)
				Expect(err).To(BeNil())
			}
			p = &rng.PersistCluster{Cluster: clust}
			node := &models.Node{ID: 1}
			chanCfg = &models.ChannelConfig{ID: uuid.New(), NodeID: node.ID}
			newRng = &models.Range{ID: uuid.New()}
			newRR = &models.RangeReplica{ID: uuid.New(), RangeID: newRng.ID, NodeID: node.ID}
			items = []interface{}{
				node,
				chanCfg,
				newRng,
				newRR,
			}
		})
		JustBeforeEach(func() {
			for _, item := range items {
				Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
			}
		})
		JustAfterEach(func() {
			for _, item := range items {
				Expect(clust.NewDelete().Model(item).WherePKs(model.NewReflect(item).PKChain().Raw()).Exec(ctx)).To(BeNil())
			}
		})
		Describe("NewRange", func() {
			It("Should save a new range, range lease, and range replica to storage", func() {
				rng, err := p.NewRange(ctx, 1)
				Expect(err).To(BeNil())
				Expect(model.NewPK(rng.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.RangeReplica.ID).IsZero()).To(BeFalse())
			})
		})
		Describe("New Range Replica", func() {
			It("Should save the replica with the correct node id", func() {
				p := &rng.PersistCluster{Cluster: clust}
				rngReplica, err := p.NewRangeReplica(ctx, newRng.ID, 1)
				Expect(err).To(BeNil())
				Expect(rngReplica.NodeID).To(Equal(1))
			})
		})
		Describe("Retrieve Range Replica", func() {
			It("Should retrieve the correct replica", func() {
				rngReplica, err := p.NewRangeReplica(ctx, newRng.ID, 1)
				Expect(err).To(BeNil())
				rr, err := p.RetrieveRangeReplicas(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(rr).To(HaveLen(2))
				Expect(rr[0].ID).To(BeElementOf(rngReplica.ID, newRR.ID))
			})
		})
		Describe("Retrieve a Range", func() {
			It("Should retrieve the correct range", func() {
				resRng, err := p.RetrieveRange(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(resRng.ID).To(Equal(newRng.ID))
			})
		})
		Describe("Retrieve Range Chunk Replicas", func() {
			It("Should retrieve the list of chunk replicas", func() {
				for i := 0; i < 10; i++ {
					cc := &models.ChannelChunk{ID: uuid.New(), RangeID: newRng.ID, ChannelConfigID: chanCfg.ID}
					Expect(clust.NewCreate().
						Model(cc).Exec(ctx)).To(BeNil())
					Expect(clust.NewCreate().
						Model(&models.ChannelChunkReplica{ID: uuid.New(), RangeReplicaID: newRR.ID, ChannelChunkID: cc.ID, Telem: telem.NewBulk([]byte("Hello"))}).Exec(ctx)).To(BeNil())
				}
				resCCR, err := p.RetrieveRangeChunkReplicas(ctx, newRng.ID)
				Expect(resCCR).To(HaveLen(10))
				Expect(err).To(BeNil())
			})
		})
		Describe("Retrieve Range Chunks", func() {
			It("Should retrieve the list of chunks", func() {
				for i := 0; i < 10; i++ {
					cc := &models.ChannelChunk{ID: uuid.New(), RangeID: newRng.ID, ChannelConfigID: chanCfg.ID}
					Expect(clust.NewCreate().
						Model(cc).Exec(ctx)).To(BeNil())
				}
				cc, err := p.RetrieveRangeChunks(ctx, newRng.ID)
				Expect(cc).To(HaveLen(10))
				Expect(err).To(BeNil())
			})
		})
		Describe("Reallocate Chunks", func() {
			It("Should reallocate the chunks correctly", func() {
				for i := 0; i < 2; i++ {
					cc := &models.ChannelChunk{ID: uuid.New(), RangeID: newRng.ID, ChannelConfigID: chanCfg.ID}
					Expect(clust.NewCreate().
						Model(cc).Exec(ctx)).To(BeNil())
				}
				cc, err := p.RetrieveRangeChunks(ctx, newRng.ID)
				Expect(err).To(BeNil())
				rng, err := p.NewRange(ctx, 1)
				Expect(err).To(BeNil())
				var ccPKs []uuid.UUID
				for _, c := range cc {
					ccPKs = append(ccPKs, c.ID)
				}
				Expect(p.ReallocateChunks(ctx, ccPKs, rng.ID)).To(BeNil())
			})
		})
		Describe("Reallocate Chunk Replicas", func() {
			It("Should retrieve the list of chunk replicas", func() {
				for i := 0; i < 10; i++ {
					cc := &models.ChannelChunk{ID: uuid.New(), RangeID: newRng.ID, ChannelConfigID: chanCfg.ID}
					Expect(clust.NewCreate().
						Model(cc).Exec(ctx)).To(BeNil())
					Expect(clust.NewCreate().
						Model(&models.ChannelChunkReplica{ID: uuid.New(), RangeReplicaID: newRR.ID, ChannelChunkID: cc.ID, Telem: telem.NewBulk([]byte("Hello"))}).Exec(ctx)).To(BeNil())
				}
				resCCR, err := p.RetrieveRangeChunkReplicas(ctx, newRng.ID)
				Expect(resCCR).To(HaveLen(10))
				Expect(err).To(BeNil())
				var ccrPKs []uuid.UUID
				for _, c := range resCCR {
					ccrPKs = append(ccrPKs, c.ID)
				}
				rng, err := p.NewRange(ctx, 1)
				cc, err := p.RetrieveRangeChunks(ctx, newRng.ID)
				Expect(err).To(BeNil())
				var ccPKs []uuid.UUID
				for _, c := range cc {
					ccPKs = append(ccPKs, c.ID)
				}
				Expect(p.ReallocateChunks(ctx, ccPKs, rng.ID)).To(BeNil())
				Expect(err).To(BeNil())
				log.Info(rng.ID, rng.RangeLease.RangeReplica.ID)
				Expect(p.ReallocateChunkReplicas(ctx, ccrPKs, rng.RangeLease.RangeReplica.ID)).To(BeNil())
				updatedResCCR, err := p.RetrieveRangeChunkReplicas(ctx, rng.ID)
				Expect(err).To(BeNil())
				Expect(updatedResCCR).To(HaveLen(10))
			})
		})
		Describe("Retrieve Range Size", func() {
			It("Should retrieve the correct range size", func() {
				for i := 0; i < 10; i++ {
					cc := &models.ChannelChunk{ID: uuid.New(), RangeID: newRng.ID, ChannelConfigID: chanCfg.ID, Size: 300}
					Expect(clust.NewCreate().
						Model(cc).Exec(ctx)).To(BeNil())
				}
				size, err := p.RetrieveRangeSize(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(size).To(Equal(int64(10 * 300)))
			})
		})
		Describe("Update Range Status", func() {
			It("Should update the range status correctly", func() {
				Expect(p.UpdateRangeStatus(ctx, newRng.ID, models.RangeStatusClosed)).To(BeNil())
				rng, err := p.RetrieveRange(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(rng.Status).To(Equal(models.RangeStatusClosed))
			})
		})
	})
})
